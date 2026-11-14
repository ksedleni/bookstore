package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/spf13/cobra"

	"github.com/ksedleni/bookstore/content"
	bookstoretypes "github.com/ksedleni/bookstore/x/bookstore/types"
)

const contentDirFlag = "content-dir"

type keyMaterial struct {
	Name string
	Addr string
	Priv *btcec.PrivateKey
	Pub  *btcec.PublicKey
}

// ContentCmd is the portable off-chain content delivery CLI.
func ContentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content",
		Short: "Owner-gated off-chain content put/get/rewrap (SHA-256 verified)",
		Long: `Publish and fetch BookStore file bytes off-chain.

The chain stores only the content hash. Envelopes are AES-256-GCM encrypted with
a CEK wrapped to the owner's secp256k1 pubkey (scheme secp256k1-ecdh-hkdf-v1).
Any third-party client can implement the same envelope format.

After a P2P sale, run "content rewrap" so the buyer can decrypt (bookstore-less;
no protocol fee).`,
	}
	cmd.PersistentFlags().String(contentDirFlag, "", "content store directory (default: <home>/content)")
	cmd.AddCommand(contentPutCmd(), contentGetCmd(), contentRewrapCmd(), contentShowCmd())
	return cmd
}

func contentStoreFromCmd(cmd *cobra.Command) (*content.Store, client.Context, error) {
	clientCtx := client.GetClientContextFromCmd(cmd)
	dir, _ := cmd.Flags().GetString(contentDirFlag)
	if dir == "" {
		dir = filepath.Join(clientCtx.HomeDir, "content")
	}
	st, err := content.NewStore(dir)
	return st, clientCtx, err
}

func loadKey(cmd *cobra.Command, name, passphrase string) (*keyMaterial, error) {
	clientCtx := client.GetClientContextFromCmd(cmd)
	rec, err := clientCtx.Keyring.Key(name)
	if err != nil {
		return nil, err
	}
	addr, err := rec.GetAddress()
	if err != nil {
		return nil, err
	}
	armored, err := clientCtx.Keyring.ExportPrivKeyArmor(name, passphrase)
	if err != nil {
		return nil, fmt.Errorf("export privkey: %w (use --keyring-backend test for localnet)", err)
	}
	priv, _, err := sdkcrypto.UnarmorDecryptPrivKey(armored, passphrase)
	if err != nil {
		return nil, err
	}
	btcecPriv, err := content.PrivFromCosmos(priv)
	if err != nil {
		return nil, err
	}
	return &keyMaterial{
		Name: name,
		Addr: addr.String(),
		Priv: btcecPriv,
		Pub:  btcecPriv.PubKey(),
	}, nil
}

func loadPubFromKeyring(cmd *cobra.Command, name string) (string, *btcec.PublicKey, error) {
	clientCtx := client.GetClientContextFromCmd(cmd)
	rec, err := clientCtx.Keyring.Key(name)
	if err != nil {
		return "", nil, err
	}
	addr, err := rec.GetAddress()
	if err != nil {
		return "", nil, err
	}
	pub, err := rec.GetPubKey()
	if err != nil {
		return "", nil, err
	}
	btcecPub, err := content.PubFromCosmos(pub)
	if err != nil {
		return "", nil, err
	}
	return addr.String(), btcecPub, nil
}

func contentPutCmd() *cobra.Command {
	var (
		from       string
		filePath   string
		originalID string
		passphrase string
	)
	cmd := &cobra.Command{
		Use:   "put",
		Short: "Encrypt and store file bytes under their SHA-256 content hash",
		RunE: func(cmd *cobra.Command, args []string) error {
			st, _, err := contentStoreFromCmd(cmd)
			if err != nil {
				return err
			}
			km, err := loadKey(cmd, from, passphrase)
			if err != nil {
				return err
			}
			plain, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			hash := content.HashPlaintext(plain)
			env, err := content.Encrypt(plain, km.Addr, km.Pub)
			if err != nil {
				return err
			}
			env.OriginalID = originalID
			if err := st.Put(env); err != nil {
				return err
			}
			out := map[string]any{
				"content_hash": hash,
				"store_path":   filepath.Join(st.Root, hash+".json"),
				"recipient":    km.Addr,
				"original_id":  originalID,
				"note":         "mint-original must use this content_hash; chain stores digest only",
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(out)
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "key name (owner/publisher)")
	cmd.Flags().StringVar(&filePath, "file", "", "path to plaintext file")
	cmd.Flags().StringVar(&originalID, "original-id", "", "optional original id for local bookkeeping")
	cmd.Flags().StringVar(&passphrase, "passphrase", "", "keyring passphrase (empty for test backend)")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("file")
	flags.AddKeyringFlags(cmd.Flags())
	return cmd
}

func contentGetCmd() *cobra.Command {
	var (
		from        string
		contentHash string
		outPath     string
		passphrase  string
		requireOwn  bool
		originalID  string
	)
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Decrypt stored content as owner and verify SHA-256",
		RunE: func(cmd *cobra.Command, args []string) error {
			st, clientCtx, err := contentStoreFromCmd(cmd)
			if err != nil {
				return err
			}
			km, err := loadKey(cmd, from, passphrase)
			if err != nil {
				return err
			}
			if contentHash == "" && originalID != "" {
				hash, owner, qerr := queryOriginalContentHash(cmd, clientCtx, originalID)
				if qerr != nil {
					return qerr
				}
				contentHash = hash
				if requireOwn && owner != km.Addr {
					return fmt.Errorf("ledger owner is %s, not %s (ownership check)", owner, km.Addr)
				}
			}
			if contentHash == "" {
				return fmt.Errorf("--content-hash or --original-id required")
			}
			env, err := st.Get(contentHash)
			if err != nil {
				return err
			}
			if requireOwn {
				if _, ok := env.FindWrap(km.Addr); !ok {
					return fmt.Errorf("no CEK wrap for %s (run content rewrap after sale?)", km.Addr)
				}
			}
			plain, err := content.Decrypt(env, km.Addr, km.Priv)
			if err != nil {
				return err
			}
			if outPath == "" {
				outPath = contentHash + ".bin"
			}
			if err := os.WriteFile(outPath, plain, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s (%d bytes) hash=%s verified\n", outPath, len(plain), env.ContentHash)
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "key name (must hold CEK wrap)")
	cmd.Flags().StringVar(&contentHash, "content-hash", "", "SHA-256 hex of plaintext")
	cmd.Flags().StringVar(&originalID, "original-id", "", "resolve content-hash from chain original id")
	cmd.Flags().StringVar(&outPath, "out", "", "output file path")
	cmd.Flags().StringVar(&passphrase, "passphrase", "", "keyring passphrase")
	cmd.Flags().BoolVar(&requireOwn, "require-owner", true, "require CEK wrap / ledger owner match when original-id set")
	_ = cmd.MarkFlagRequired("from")
	flags.AddKeyringFlags(cmd.Flags())
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func contentRewrapCmd() *cobra.Command {
	var (
		from        string
		to          string
		contentHash string
		passphrase  string
	)
	cmd := &cobra.Command{
		Use:   "rewrap",
		Short: "Re-wrap CEK for a new owner after P2P sale/transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			st, _, err := contentStoreFromCmd(cmd)
			if err != nil {
				return err
			}
			km, err := loadKey(cmd, from, passphrase)
			if err != nil {
				return err
			}
			toAddr, toPub, err := loadPubFromKeyring(cmd, to)
			if err != nil {
				return err
			}
			env, err := st.Get(contentHash)
			if err != nil {
				return err
			}
			if err := content.Rewrap(env, km.Addr, km.Priv, toAddr, toPub); err != nil {
				return err
			}
			if err := st.Put(env); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "rewrapped %s for %s (%s)\n", contentHash, to, toAddr)
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "current key that can unwrap CEK")
	cmd.Flags().StringVar(&to, "to", "", "new owner key name")
	cmd.Flags().StringVar(&contentHash, "content-hash", "", "SHA-256 hex")
	cmd.Flags().StringVar(&passphrase, "passphrase", "", "keyring passphrase")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")
	_ = cmd.MarkFlagRequired("content-hash")
	flags.AddKeyringFlags(cmd.Flags())
	return cmd
}

func contentShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [content-hash]",
		Short: "Show envelope metadata (no decrypt)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, _, err := contentStoreFromCmd(cmd)
			if err != nil {
				return err
			}
			env, err := st.Get(args[0])
			if err != nil {
				return err
			}
			recipients := make([]string, 0, len(env.KeyWraps))
			for _, w := range env.KeyWraps {
				recipients = append(recipients, w.Recipient)
			}
			summary := map[string]any{
				"version":          env.Version,
				"content_hash":     env.ContentHash,
				"cipher":           env.Cipher,
				"wraps":            len(env.KeyWraps),
				"recipients":       recipients,
				"original_id":      env.OriginalID,
				"ciphertext_bytes": len(env.Ciphertext) / 2,
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(summary)
		},
	}
	return cmd
}

func queryOriginalContentHash(cmd *cobra.Command, clientCtx client.Context, originalID string) (string, string, error) {
	id, err := strconv.ParseUint(originalID, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("invalid original-id: %w", err)
	}
	clientCtx, err = client.GetClientQueryContext(cmd)
	if err != nil {
		return "", "", err
	}
	qc := bookstoretypes.NewQueryClient(clientCtx)
	res, err := qc.Original(cmd.Context(), &bookstoretypes.QueryGetOriginalRequest{Id: id})
	if err != nil {
		return "", "", err
	}
	return res.Original.ContentHash, res.Original.Owner, nil
}
