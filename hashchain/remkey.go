package hashchain

import (
	"fmt"

	"github.com/frankbraun/codechain/internal/base64"
	"github.com/frankbraun/codechain/util/time"
)

// RemoveKey adds a pubkey remove entry to hash chain.
func (c *HashChain) RemoveKey(pubKey [32]byte) (string, error) {
	// check arguments
	// TODO: check that pubkey is actually active in chain
	// TODO: check that still enough public keys remain to reach m

	// create entry
	prev := c.LastEntryHash()
	l := &link{
		previous:   prev[:],
		datum:      time.Now(),
		linkType:   removeKeyType,
		typeFields: []string{base64.Encode(pubKey[:])},
	}
	c.chain = append(c.chain, l)
	entry := l.String()

	// verify
	if err := c.verify(); err != nil {
		return "", err
	}

	// save
	if _, err := fmt.Fprintln(c.fp, entry); err != nil {
		return "", err
	}
	return entry, nil
}
