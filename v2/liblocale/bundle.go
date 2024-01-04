package liblocale

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/helloferdie/golib/v2/liblogger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// LoadBundle
func LoadBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(detectDefaultLocale())

	// Load localization
	dirLocale := os.Getenv("dir_locale")
	files, err := os.ReadDir(dirLocale)
	if err != nil {
		liblogger.Errorf("Error read localization directory %v", err)
		return nil, err
	}
	for _, file := range files {
		// Process localization directory
		if file.IsDir() {
			languageName := file.Name()
			_, err := language.Parse(languageName) // Check for valid language
			if err != nil {
				liblogger.Errorf("Error localization directory %s is not a valid language tag", languageName)
				continue
			}

			dirLang := dirLocale + "/" + languageName
			languageFiles, err := os.ReadDir(dirLang)
			for _, languageFile := range languageFiles {
				if !languageFile.IsDir() {
					processLanguageFile(bundle, languageName, dirLang, languageFile)
				}
			}
		}
	}
	return bundle, nil
}

// detectDefaultLocale detects the default locale from environment
func detectDefaultLocale() language.Tag {
	defaultLocale := os.Getenv("locale_default")
	if defaultLocale == "" {
		return language.English
	}
	tag, err := language.Parse(defaultLocale)
	if err != nil {
		liblogger.Errorf("Error parsing default locale: %v", err)
		return language.English
	}
	return tag
}

// processLanguageFile read language file and load to bundle
func processLanguageFile(bundle *i18n.Bundle, languageName string, languageDir string, languageFile fs.DirEntry) error {
	filename := languageFile.Name()
	extension := path.Ext(filename)
	if extension == ".json" {
		languageFilename := languageDir + "/" + filename
		b, err := os.ReadFile(languageFilename)
		if err != nil {
			liblogger.Errorf("Error localization failed to read %s", languageFilename)
			return err
		}

		data := map[string]interface{}{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			liblogger.Errorf("Error localization failed to parse %s", languageFilename)
			return err
		}

		// Flatten json
		flattenData := make(map[string]string)
		flattenJSON(strings.TrimSuffix(filename, extension), data, flattenData)
		flattenBytes, err := json.Marshal(flattenData)
		if err != nil {
			liblogger.Errorf("Error localization failed to flatten %s", languageFilename)
			return err
		}

		// Parse and set to bundle
		_, err = bundle.ParseMessageFileBytes(flattenBytes, languageName+extension)
		if err != nil {
			liblogger.Errorf("Error localization bundle failed to parse %s", languageFilename)
			return err
		}
	}
	return nil
}

// flattenJSON convert map[string]interface{}{} to map[string]string, where the key delimiter with "."
func flattenJSON(prefix string, value interface{}, out map[string]string) {
	// Check the type of the value
	switch value := value.(type) {
	case map[string]interface{}:
		// If it's a map, iterate through its keys
		for k, v := range value {
			// Construct the new key
			newKey := k
			if prefix != "" {
				newKey = prefix + "." + k
			}
			// Recursively flatten the value
			flattenJSON(newKey, v, out)
		}
	case string:
		// If it's a string, add it to the output
		out[prefix] = value
	}
}
