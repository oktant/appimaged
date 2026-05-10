package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("== appimaged ==")
	fmt.Println()

	fmt.Print("Path to AppImage: ")
	appImagePath, _ := reader.ReadString('\n')
	appImagePath = strings.TrimSpace(appImagePath)

	if _, err := os.Stat(appImagePath); err != nil {
		fmt.Println("AppImage not found:", err)
		os.Exit(1)
	}

	fmt.Print("Application name: ")
	appName, _ := reader.ReadString('\n')
	appName = strings.TrimSpace(appName)

	fmt.Print("Path to icon (.png/.svg) [optional]: ")
	iconPath, _ := reader.ReadString('\n')
	iconPath = strings.TrimSpace(iconPath)

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home directory:", err)
		os.Exit(1)
	}

	appDir := getApplicationsDir(home)
	iconDir := filepath.Join(home, ".local", "share", "icons")
	desktopDir := filepath.Join(home, ".local", "share", "applications")

	createDir(appDir)
	createDir(iconDir)
	createDir(desktopDir)

	// Copy AppImage
	appImageName := filepath.Base(appImagePath)
	destAppImage := filepath.Join(appDir, appImageName)

	err = copyFile(appImagePath, destAppImage)
	if err != nil {
		fmt.Println("Failed to copy AppImage:", err)
		os.Exit(1)
	}

	// Make executable
	err = os.Chmod(destAppImage, 0755)
	if err != nil {
		fmt.Println("Failed to make AppImage executable:", err)
		os.Exit(1)
	}

	// Copy icon
	iconDest := ""
	if iconPath != "" {
		iconName := filepath.Base(iconPath)
		iconDest = filepath.Join(iconDir, iconName)

		err = copyFile(iconPath, iconDest)
		if err != nil {
			fmt.Println("Failed to copy icon:", err)
			os.Exit(1)
		}
	}

	// Create desktop entry
	desktopFileName := sanitize(appName) + ".desktop"
	desktopFilePath := filepath.Join(desktopDir, desktopFileName)

	desktopContent := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=%s
Exec=%s
Icon=%s
Terminal=false
Categories=Utility;
`, appName, destAppImage, iconDest)

	err = os.WriteFile(desktopFilePath, []byte(desktopContent), 0755)
	if err != nil {
		fmt.Println("Failed to create desktop entry:", err)
		os.Exit(1)
	}

	// Refresh desktop database (optional)
	exec.Command("update-desktop-database", desktopDir).Run()

	fmt.Println()
	fmt.Println("AppImage installed successfully!")
	fmt.Println("AppImage:", destAppImage)
	fmt.Println("Launcher:", desktopFilePath)
	fmt.Println("Search for the app in your desktop dashboard/menu.")
}

func getApplicationsDir(home string) string {
	systemApplications := "/Applications"

	// Prefer /Applications if it exists
	info, err := os.Stat(systemApplications)
	if err == nil && info.IsDir() {
		return systemApplications
	}

	// Fallback to ~/Applications
	return filepath.Join(home, "Applications")
}

func createDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", path, err)
		os.Exit(1)
	}
}

func sanitize(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func copyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = output.ReadFrom(input)
	if err != nil {
		return err
	}

	return output.Sync()
}