package filepreview

import (
	"encoding/base64"
	"os"
)

func pathToBase64Encode(filePath string) (string, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return "", err
    }

    encodedData := base64.StdEncoding.EncodeToString(data)
    return encodedData, nil
}
