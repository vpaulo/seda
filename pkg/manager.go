package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manager handles package operations
type Manager struct {
	PackagesDir string
}

// NewManager creates a new package manager
func NewManager() *Manager {
	home_dir, _ := os.UserHomeDir()
	packages_dir := filepath.Join(home_dir, ".seda", "packages")

	// Ensure packages directory exists
	os.MkdirAll(packages_dir, 0755)

	return &Manager{
		PackagesDir: packages_dir,
	}
}

// Install downloads and caches a package from a git repository
func (m *Manager) Install(repo_url string) error {
	// Parse repository URL to get package name
	package_name := m.get_package_name(repo_url)
	package_path := filepath.Join(m.PackagesDir, package_name)

	// Check if package already exists
	if _, err := os.Stat(package_path); err == nil {
		fmt.Printf("Package %s already installed\n", package_name)
		return nil
	}

	fmt.Printf("Installing package %s...\n", package_name)

	// Ensure parent directories exist
	parent_dir := filepath.Dir(package_path)
	if err := os.MkdirAll(parent_dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %v", err)
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", repo_url, package_path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	fmt.Printf("Package %s installed successfully\n", package_name)
	return nil
}

// Update updates an existing package
func (m *Manager) Update(package_name string) error {
	package_path := filepath.Join(m.PackagesDir, package_name)

	// Check if package exists
	if _, err := os.Stat(package_path); os.IsNotExist(err) {
		return fmt.Errorf("package %s not found", package_name)
	}

	fmt.Printf("Updating package %s...\n", package_name)

	// Pull latest changes
	cmd := exec.Command("git", "pull")
	cmd.Dir = package_path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update package: %v", err)
	}

	fmt.Printf("Package %s updated successfully\n", package_name)
	return nil
}

// Remove removes a package from the cache
func (m *Manager) Remove(package_name string) error {
	package_path := filepath.Join(m.PackagesDir, package_name)

	// Check if package exists
	if _, err := os.Stat(package_path); os.IsNotExist(err) {
		return fmt.Errorf("package %s not found", package_name)
	}

	fmt.Printf("Removing package %s...\n", package_name)

	if err := os.RemoveAll(package_path); err != nil {
		return fmt.Errorf("failed to remove package: %v", err)
	}

	fmt.Printf("Package %s removed successfully\n", package_name)
	return nil
}

// List lists all installed packages
func (m *Manager) List() error {
	packages := []string{}

	err := filepath.Walk(m.PackagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == m.PackagesDir {
			return nil
		}

		if !info.IsDir() && info.Name() == "module.s" {
			packages = append(packages, filepath.Dir(path)[len(m.PackagesDir):]+"/"+info.Name())
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to read packages directory: %v", err)
	}

	if len(packages) == 0 {
		fmt.Println("No packages installed")
		return nil
	}

	fmt.Println("Installed packages:")
	for _, pkg := range packages {
		fmt.Printf("  - %s\n", pkg[1:])
	}

	return nil
}

// GetPackagePath returns the full path to a package
func (m *Manager) GetPackagePath(package_name string) string {
	return filepath.Join(m.PackagesDir, package_name)
}

// IsInstalled checks if a package is installed
func (m *Manager) IsInstalled(package_name string) bool {
	package_path := filepath.Join(m.PackagesDir, package_name)
	_, err := os.Stat(package_path)
	return err == nil
}

// get_package_name extracts package name from repository URL
func (m *Manager) get_package_name(repo_url string) string {
	// Handle different URL formats
	if strings.HasPrefix(repo_url, "https://github.com/") {
		// Extract user/repo from GitHub URL
		// Keep directory structure: github.com/user/repo
		parts := strings.TrimPrefix(repo_url, "https://")
		parts = strings.TrimSuffix(parts, ".git")
		return parts
	}

	if strings.HasPrefix(repo_url, "https://gitlab.com/") {
		// Extract user/repo from GitLab URL
		parts := strings.TrimPrefix(repo_url, "https://")
		parts = strings.TrimSuffix(parts, ".git")
		return parts
	}

	if strings.Contains(repo_url, "/") {
		// For URLs like github.com/user/repo
		parts := strings.TrimSuffix(repo_url, ".git")
		return parts
	}

	// Default: use the last part of the URL
	parts := strings.Split(repo_url, "/")
	name := parts[len(parts)-1]
	return strings.TrimSuffix(name, ".git")
}
