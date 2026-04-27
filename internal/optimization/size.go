package optimization

import (
    "debug/pe"
    "os"
    "os/exec"
)

func OptimizeBinary(path string) error {
    fi, err := os.Stat(path)
    if err != nil {
        return err
    }

    originalSize := fi.Size()

    if err := stripDebugInfo(path); err != nil {
        return err
    }

    if err := optimizeWithUPX(path); err != nil {
        return err
    }

    newFi, err := os.Stat(path)
    if err != nil {
        return err
    }

    _ = originalSize - newFi.Size()

    return nil
}

func stripDebugInfo(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return nil
    }
    defer f.Close()

    peFile, err := pe.NewFile(f)
    if err != nil {
        return nil
    }

    hasDebugSections := false
    for _, section := range peFile.Sections {
        name := section.Name
        if name == ".debug_info" || name == ".debug_line" || name == ".debug_abbrev" {
            hasDebugSections = true
            break
        }
    }

    if !hasDebugSections {
        return nil
    }

    _, err = exec.LookPath("strip")
    if err != nil {
        return nil
    }

    cmd := exec.Command("strip", "--strip-debug", path)
    if output := cmd.Run(); output != nil {
        return nil
    }

    return nil
}

func trimNulBytes_(b []byte) string {
    for i, c := range b {
        if c == 0 {
            return string(b[:i])
        }
    }
    return string(b)
}

func optimizeWithUPX(path string) error {
    _, err := exec.LookPath("upx")
    if err != nil {
        return nil
    }

    cmd := exec.Command("upx", "-9", path)
    if output := cmd.Run(); output != nil {
        return nil
    }

    return nil
}

func GetBinarySize(path string) (int64, error) {
    fi, err := os.Stat(path)
    if err != nil {
        return 0, err
    }
    return fi.Size(), nil
}