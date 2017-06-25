package gogurt

import (
	// "log"
	"os/exec"
)

// Just call sed.
func ReplaceInFile(fileName, pattern, repl string) error {
	cmd := exec.Command(
		"sed",
		"--in-place",
		"--regexp-extended",
		"s/" + pattern + "/" + repl + "/g",
		fileName,
	)
	// log.Println("Sed command: ", cmd)

	return cmd.Run()

	// patt, err := regexp.Compile("\n" + pattern + "\n")
	// if err != nil {
	// 	return err
	// }

	// stat, err := os.Stat(fileName)
	// if err != nil {
	// 	return err
	// }

	// fileContents, err := ioutil.ReadFile(fileName)
	// if err != nil {
	// 	return err
	// }

	// replacedContents := patt.ReplaceAll(fileContents, []byte(repl))
	// ioutil.WriteFile(fileName, replacedContents, stat.Mode())

	// return nil
}
