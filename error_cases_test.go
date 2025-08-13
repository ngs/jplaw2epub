package jplaw2epub

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ngs.io/jplaw-xml"
)

func TestCreateEPUBFromXMLFile_Errors(t *testing.T) {
	tests := []struct {
		name    string
		xmlData string
		wantErr bool
	}{
		{
			name:    "Invalid XML",
			xmlData: "not valid xml",
			wantErr: true,
		},
		{
			name:    "Empty XML",
			xmlData: "",
			wantErr: true,
		},
		{
			name: "XML without LawBody",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
				<Law Era="Reiwa" Year="1" Num="1" LawType="Act" Lang="ja">
					<LawNum>テスト</LawNum>
				</Law>`,
			wantErr: true,
		},
		{
			name: "XML without LawTitle",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
				<Law Era="Reiwa" Year="1" Num="1" LawType="Act" Lang="ja">
					<LawNum>テスト</LawNum>
					<LawBody>
						<MainProvision></MainProvision>
					</LawBody>
				</Law>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpfile, err := os.CreateTemp("", "error_test_*.xml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, writeErr := tmpfile.WriteString(tt.xmlData); writeErr != nil {
				t.Fatalf("Failed to write test XML: %v", writeErr)
			}
			if closeErr := tmpfile.Close(); closeErr != nil {
				t.Fatalf("Failed to close temp file: %v", closeErr)
			}

			// Test
			xmlFile, err := os.Open(tmpfile.Name())
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			defer xmlFile.Close()

			_, err = CreateEPUBFromXMLFile(xmlFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEPUBFromXMLFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateEPUBFromXMLPath_Errors(t *testing.T) {
	// Test with non-existent file
	_, err := CreateEPUBFromXMLPath("/non/existent/file.xml")
	if err == nil {
		t.Error("CreateEPUBFromXMLPath should return error for non-existent file")
	}
}

func TestWriteEPUB_Errors(t *testing.T) {
	// Create a simple EPUB
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<Law Era="Reiwa" Year="1" Num="1" LawType="Act" Lang="ja">
  <LawNum>テスト</LawNum>
  <LawBody>
    <LawTitle>テスト法</LawTitle>
    <MainProvision>
      <Chapter Num="1">
        <ChapterTitle>第一章</ChapterTitle>
      </Chapter>
    </MainProvision>
  </LawBody>
</Law>`

	tmpfile, err := os.CreateTemp("", "test_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, writeErr := tmpfile.WriteString(xmlContent); writeErr != nil {
		t.Fatalf("Failed to write test XML: %v", writeErr)
	}
	if closeErr := tmpfile.Close(); closeErr != nil {
		t.Fatalf("Failed to close temp file: %v", closeErr)
	}

	book, err := CreateEPUBFromXMLPath(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	// Test writing to invalid path
	err = WriteEPUB(book, "/invalid\x00path/test.epub") // Path with null character
	if err == nil {
		t.Error("WriteEPUB should return error for invalid path")
	}
}

func TestProcessChapterWithImages_Error(t *testing.T) {
	book, _ := createEPUBFromData(&jplaw.Law{
		LawBody: jplaw.LawBody{
			LawTitle: &jplaw.LawTitle{Content: "テスト"},
		},
	})

	// Create chapter with article that has very long title (potential issue)
	longTitle := strings.Repeat("あ", 10000)
	chapter := &jplaw.Chapter{
		ChapterTitle: jplaw.ChapterTitle{
			Content: longTitle,
		},
	}

	// This should not panic
	err := processChapterWithImages(book, chapter, 0, nil)
	if err != nil {
		// Error is acceptable, panic is not
		t.Logf("processChapterWithImages returned error (expected): %v", err)
	}
}

func TestLoadXMLDataFromReader_InvalidXML(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "Invalid XML structure",
			data:    "<unclosed",
			wantErr: true,
		},
		{
			name:    "Non-XML data",
			data:    "This is not XML",
			wantErr: true,
		},
		{
			name:    "Empty data",
			data:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.data)
			_, err := loadXMLDataFromReader(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadXMLDataFromReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateEPUBFromData_MinimalLaw(t *testing.T) {
	// Test with minimal valid Law structure
	data := &jplaw.Law{
		LawBody: jplaw.LawBody{
			LawTitle: &jplaw.LawTitle{
				Content: "最小限の法律",
			},
			MainProvision: jplaw.MainProvision{
				// No chapters, articles, or paragraphs
			},
		},
	}

	book, err := createEPUBFromData(data)
	if err != nil {
		t.Errorf("createEPUBFromData() with minimal law should not error: %v", err)
	}
	if book == nil {
		t.Error("createEPUBFromData() should return non-nil book for minimal law")
	}
}

func TestWriteEPUB_DirectoryCreation(t *testing.T) {
	// Create a simple EPUB
	book, err := createEPUBFromData(&jplaw.Law{
		LawBody: jplaw.LawBody{
			LawTitle: &jplaw.LawTitle{Content: "テスト"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	// Test writing to nested non-existent directory
	tmpDir := t.TempDir()
	nestedPath := filepath.Join(tmpDir, "nested", "dir", "test.epub")

	err = WriteEPUB(book, nestedPath)
	if err != nil {
		t.Errorf("WriteEPUB should create nested directories: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
		t.Error("WriteEPUB did not create the file in nested directory")
	}
}
