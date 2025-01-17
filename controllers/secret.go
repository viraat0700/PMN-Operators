// /*
// Copyright 2024.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

// package controllers

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"path/filepath"

// 	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"
// 	corev1 "k8s.io/api/core/v1"
// 	apierrors "k8s.io/apimachinery/pkg/api/errors"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	ctrl "sigs.k8s.io/controller-runtime"
// )

// func (r *PmnsystemReconciler) CreateSecretsFromCertificates(secretName string, certDir string, requiredFiles []string, namespace string, cr *v1.Pmnsystem) error {
// 	// Logger for debugging
// 	log := ctrl.Log.WithName("CreateSecretsFromCertificates")
// 	log.Info("Starting secret creation", "SecretName", secretName, "Namespace", namespace, "CertificateDir", certDir)

// 	// 1. Log Current Working Directory
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		log.Error(err, "Failed to get current working directory")
// 		return fmt.Errorf("failed to get current working directory: %w", err)
// 	}
// 	log.Info("Current working directory", "WorkingDir", cwd)

// 	// 2. Resolve Absolute Path for the Certificate Directory
// 	absolutePath, err := filepath.Abs(certDir)
// 	if err != nil {
// 		log.Error(err, "Failed to resolve absolute path for certDir", "certDir", certDir)
// 		return fmt.Errorf("failed to resolve absolute path for %s: %w", certDir, err)
// 	}
// 	log.Info("Resolved absolute path for certificate directory", "AbsolutePath", absolutePath)

// 	// 3. Verify Directory Existence
// 	if _, err := os.Stat(certDir); os.IsNotExist(err) {
// 		log.Error(err, "Certificate directory does not exist", "certDir", certDir)
// 		return fmt.Errorf("certificate directory %s does not exist", certDir)
// 	} else if err != nil {
// 		log.Error(err, "Error accessing certificate directory", "certDir", certDir)
// 		return fmt.Errorf("error accessing certificate directory %s: %w", certDir, err)
// 	}
// 	log.Info("Certificate directory verified", "certDir", certDir)

// 	// 4. List All Files in the Directory for Debugging
// 	files, err := os.ReadDir(certDir)
// 	if err != nil {
// 		log.Error(err, "Failed to read directory", "Directory", certDir)
// 		return fmt.Errorf("failed to read directory %s: %w", certDir, err)
// 	}
// 	log.Info("Files in certificate directory:")
// 	for _, f := range files {
// 		log.Info("File found", "FileName", f.Name())
// 	}

// 	// 5. Validate and Read Required Files
// 	data := make(map[string][]byte)
// 	for _, file := range requiredFiles {
// 		// Construct the full file path
// 		filePath := filepath.Join(certDir, file)
// 		log.Info("Checking required file", "FileName", file, "FilePath", filePath)

// 		// Verify file existence
// 		if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 			log.Error(err, "Required file not found", "FilePath", filePath)
// 			return fmt.Errorf("required file %s is missing in directory %s", file, certDir)
// 		} else if err != nil {
// 			log.Error(err, "Error accessing required file", "FilePath", filePath)
// 			return fmt.Errorf("error accessing file %s in directory %s: %w", file, certDir, err)
// 		}

// 		// Read file content
// 		content, err := os.ReadFile(filePath)
// 		if err != nil {
// 			log.Error(err, "Failed to read file", "FilePath", filePath)
// 			return fmt.Errorf("failed to read file %s: %w", filePath, err)
// 		}
// 		data[file] = content
// 		log.Info("File read successfully", "FileName", file, "FilePath", filePath)
// 	}

// 	// 6. Create Kubernetes Secret
// 	secret := &corev1.Secret{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      secretName,
// 			Namespace: namespace,
// 			OwnerReferences: []metav1.OwnerReference{
// 				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
// 					Group:   v1.GroupVersion.Group,
// 					Version: v1.GroupVersion.Version,
// 					Kind:    "Pmnsystem",
// 				}),
// 			},
// 		},
// 		Data: data,
// 		Type: corev1.SecretTypeOpaque,
// 	}

// 	// Attempt to Create the Secret
// 	err = r.Client.Create(context.TODO(), secret)
// 	if err != nil {
// 		if apierrors.IsAlreadyExists(err) {
// 			log.Info("Secret already exists, updating it", "SecretName", secretName)
// 			// Update the existing secret
// 			err = r.Client.Update(context.TODO(), secret)
// 			if err != nil {
// 				log.Error(err, "Failed to update existing secret", "SecretName", secretName)
// 				return fmt.Errorf("failed to update secret %s: %w", secretName, err)
// 			}
// 		} else {
// 			log.Error(err, "Failed to create secret", "SecretName", secretName)
// 			return fmt.Errorf("failed to create secret %s: %w", secretName, err)
// 		}
// 	}

//		log.Info("Secret created or updated successfully", "SecretName", secretName, "Namespace", namespace)
//		return nil
//	}
package controllers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *PmnsystemReconciler) CreateSecretsFromCertificates(secretName string, certDir string, requiredFiles []string, namespace string, cr *v1.Pmnsystem) error {
	// Logger for debugging
	log := ctrl.Log.WithName("CreateSecretsFromCertificates")
	log.Info("Starting secret creation", "SecretName", secretName, "Namespace", namespace, "CertificateDir", certDir)

	data := make(map[string][]byte)

	for _, file := range requiredFiles {
		filePath := filepath.Join(certDir, file)
		log.Info("Checking required file", "FileName", file, "FilePath", filePath)

		// Verify file existence
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Info("File not found, calling generateCertificate", "FileName", file)

			// Define parameters for generateCertificate
			repoPath := "/mnt/c/Users/MSI-1/Desktop/PMN-Operator-Two/pmn-operator/mTLS/"
			certMovePath := "/mnt/c/Users/MSI-1/Desktop/PMN-Operator-Two/pmn-operator/VMFolder"
			linesToReplace := map[string]map[string]string{
				"ca.cnf": {
					"commonName=		notifier-ca.operator.wavelabs.int":         "commonName=		notifier-ca.operator.wavelabs.int",
					"emailAddress=		viraat.shrivastava@veltris.com":          "emailAddress=		viraat.shrivastava@veltris.com",
					"subjectAltName = DNS:notifier-ca.operator.wavelabs.int": "subjectAltName = DNS:notifier-ca.operator.wavelabs.int",
				},
				"client.cnf": {
					"commonName=		notifier-client.operator.wavelabs.int":         "commonName=		notifier-client.operator.wavelabs.int",
					"emailAddress=		viraat.shrivastava@veltris.com":              "emailAddress=		viraat.shrivastava@veltris.com",
					"subjectAltName = DNS:notifier-client.operator.wavelabs.int": "subjectAltName = DNS:notifier-client.operator.wavelabs.int",
				},
				"server.cnf": {
					"commonName=		notifier-server.operator.wavelabs.int":         "commonName=		notifier-server.operator.wavelabs.int",
					"emailAddress=		viraat.shrivastava@veltris.com":              "emailAddress=		viraat.shrivastava@veltris.com",
					"subjectAltName = DNS:notifier-server.operator.wavelabs.int": "subjectAltName = DNS:notifier-server.operator.wavelabs.int",
				},
			}

			newFileNames := map[string]string{
				"ca.crt":     "notifier-ca.crt",
				"ca.key":     "notifier-ca.key",
				"client.key": "notifier-client.key",
				"client.crt": "notifier-client.crt",
				"server.key": "notifier-server.key",
				"server.crt": "notifier-server.crt",
			}

			certSubDir := "certs"
			// Call generateCertificate
			err := r.generateCertificate(repoPath, certMovePath, certSubDir, linesToReplace, newFileNames)
			if err != nil {
				log.Error(err, "Failed to generate certificate", "FileName", file)
				return fmt.Errorf("failed to generate certificate %s: %w", file, err)
			}

			log.Info("Certificate generated successfully", "FileName", file)
		} else if err != nil {
			log.Error(err, "Error accessing required file", "FilePath", filePath)
			return fmt.Errorf("error accessing file %s in directory %s: %w", file, certDir, err)
		}

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Error(err, "Failed to read file", "FilePath", filePath)
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}
		data[file] = content
		log.Info("File read successfully", "FileName", file, "FilePath", filePath)
	}

	// 6. Create Kubernetes Secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Data: data,
		Type: corev1.SecretTypeOpaque,
	}

	// Attempt to Create the Secret
	err := r.Client.Create(context.TODO(), secret)
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			log.Info("Secret already exists, updating it", "SecretName", secretName)
			// Update the existing secret
			err = r.Client.Update(context.TODO(), secret)
			if err != nil {
				log.Error(err, "Failed to update existing secret", "SecretName", secretName)
				return fmt.Errorf("failed to update secret %s: %w", secretName, err)
			}
		} else {
			log.Error(err, "Failed to create secret", "SecretName", secretName)
			return fmt.Errorf("failed to create secret %s: %w", secretName, err)
		}
	}

	log.Info("Secret created or updated successfully", "SecretName", secretName, "Namespace", namespace)
	return nil
}

// generateCertificate generates missing certificates
func (r *PmnsystemReconciler) generateCertificate(repoPath, certMovePath, certSubDir string, linesToReplace map[string]map[string]string, newFileNames map[string]string) error {
	log := ctrl.Log.WithName("GenerateCertificate")
	log.Info("Starting certificate generation", "RepoPath", repoPath, "CertMovePath", certMovePath, "CertSubDir", certSubDir)

	// Navigate to the certs directory
	certsPath := filepath.Join(repoPath, certSubDir)
	log.Info("Navigating to certs directory", "CertsPath", certsPath)

	// Step 1: Edit the required files
	for fileName, replacements := range linesToReplace {
		filePath := filepath.Join(certsPath, fileName)
		log.Info("Editing file", "FilePath", filePath)

		// Read the file
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Error(err, "Failed to read file for editing", "FilePath", filePath)
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		// Split content into lines
		lines := bytes.Split(content, []byte("\n"))
		for i, line := range lines {
			lineStr := string(line)
			// Check if the line matches any replacement
			if newValue, exists := replacements[lineStr]; exists {
				log.Info("Replacing line", "OldLine", lineStr, "NewValue", newValue)
				lines[i] = []byte(newValue)
			}
		}

		// Write the updated content back to the file
		updatedContent := bytes.Join(lines, []byte("\n"))
		err = os.WriteFile(filePath, updatedContent, 0644)
		if err != nil {
			log.Error(err, "Failed to write updated file", "FilePath", filePath)
			return fmt.Errorf("failed to write updated file %s: %w", filePath, err)
		}
		log.Info("File edited successfully", "FilePath", filePath)
	}

	// Step 2: Run `make cert`
	cmd := exec.Command("make", "cert")
	cmd.Dir = certsPath
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	log.Info("Executing make cert command")
	err := cmd.Run()
	if err != nil {
		log.Error(err, "Failed to run make cert", "Stdout", out.String(), "Stderr", stderr.String())
		return fmt.Errorf("failed to run make cert: %w", err)
	}
	log.Info("make cert command executed successfully", "Output", out.String())

	// Step 3: Rename and move generated files
	for oldName, newName := range newFileNames {
		src := filepath.Join(certsPath, oldName)
		dst := filepath.Join(certMovePath, newName)
		log.Info("Renaming and moving file", "OldName", oldName, "NewName", newName, "Source", src, "Destination", dst)

		err = os.Rename(src, dst)
		if err != nil {
			log.Error(err, "Failed to rename and move file", "OldName", oldName, "NewName", newName)
			return fmt.Errorf("failed to rename and move file %s to %s: %w", oldName, newName, err)
		}
		log.Info("File renamed and moved successfully", "OldName", oldName, "NewName", newName)
	}

	log.Info("Certificate generation completed successfully")
	return nil
}
