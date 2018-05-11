package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/hashchain"
	"github.com/frankbraun/codechain/internal/def"
	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/log"
)

// Start implements the 'start' command.
func Start(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -s seckey.bin\n", argv0)
		fmt.Fprintf(os.Stderr, "Initialized new .codechain/hashchain in current directory.\n")
		fs.PrintDefaults()
	}
	seckey := fs.String("s", "", "Secret key file")
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	if *seckey == "" {
		return fmt.Errorf("%s: option -s is mandatory", argv0)
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	if err := os.MkdirAll(def.CodechainDir, 0755); err != nil {
		return err
	}
	exists, err := file.Exists(hashchainFile)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s: file '%s' exists already", argv0, hashchainFile)
	}
	sec, _, comment, err := seckeyRead(*seckey)
	if err != nil {
		return err
	}
	c, entry, err := hashchain.Start(hashchainFile, *sec, comment)
	if err != nil {
		return err
	}
	defer c.Close()
	fmt.Println(entry)
	return nil
}
