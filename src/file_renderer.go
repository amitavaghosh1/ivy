package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

func HandleEnv(res map[string][]ConfigParameter, env string) error {
	configForEnv, ok := res[env]
	if !ok {
		return errors.New("invalid_env")
	}

	var w = bytes.Buffer{}

	for _, cfp := range configForEnv {
		w.WriteString(cfp.Name)
		w.WriteString("=")
		w.WriteString(cfp.Value)
		w.WriteString("\n")
	}

	file, err := ioutil.TempFile(os.TempDir(), "*.env")
	if err != nil {
		log.Fatal("failed to create temp file", err)
	}
	defer os.Remove(file.Name())

	err = OpenInEditor(w.Bytes(), file)
	if err != nil {
		log.Fatal("failed to write to file ", err)
	}

	return nil
}

func OpenInEditor(b []byte, file *os.File) error {
	_, err := file.Write(b)
	if err != nil {
		log.Fatal("failed to write to file ", err)
	}

	cmd := exec.Command("vim", file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Println("error closing vim ", err)
		return err
	}

	return nil
}

func HandleYAML(res map[string]map[string][]ConfigParameter) error {
	b, err := yaml.Marshal(res)
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile(os.TempDir(), "*.yaml")
	if err != nil {
		log.Fatal("failed to create temp file", err)
	}
	defer os.Remove(file.Name())

	err = OpenInEditor(b, file)
	if err != nil {
		log.Fatal("failed to write to file ", err)
	}

	return nil
}

func HandleJSON(res map[string]map[string][]ConfigParameter) error {
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile(os.TempDir(), "*.json")
	if err != nil {
		log.Fatal("failed to create temp file", err)
	}
	defer os.Remove(file.Name())

	err = OpenInEditor(b, file)
	if err != nil {
		log.Fatal("failed to write to file ", err)
	}

	return nil
}
