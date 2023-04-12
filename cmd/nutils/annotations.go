package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/notaryproject/notation-core-go/signature/cose"
	"github.com/spf13/cobra"
)

const annotationX509ChainThumbprint = "io.cncf.notary.x509chain.thumbprint#S256"

func annotationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "annotations <cose_signature>",
		Aliases: []string{"a"},
		Short:   "Prints the annotations of a COSE signature",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotations(cmd.Context(), args[0])
		},
	}
	return cmd
}

// runAnnotations prints the annotations of a COSE signature
// ref: https://github.com/notaryproject/notation-go/blob/753b6b12cb2ebba56df6bfdc5b75748ad9e0a5bd/notation.go#L404
func runAnnotations(ctx context.Context, path string) error {
	// read cose signature
	envBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	envelope, err := cose.ParseEnvelope(envBytes)
	if err != nil {
		return err
	}
	content, err := envelope.Content()
	if err != nil {
		return err
	}

	// generate annotations
	annotations := make(map[string]string)

	// generate thumbprints
	var thumbprints []string
	for _, cert := range content.SignerInfo.CertificateChain {
		thumbprint := sha256.Sum256(cert.Raw)
		thumbprints = append(thumbprints, hex.EncodeToString(thumbprint[:]))
	}
	thumbprintsJSON, err := json.Marshal(thumbprints)
	if err != nil {
		return err
	}
	annotations[annotationX509ChainThumbprint] = string(thumbprintsJSON)

	// print annotations
	return json.NewEncoder(os.Stdout).Encode(annotations)
}
