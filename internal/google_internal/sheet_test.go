package google_internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractSheetIdFromUrl(t *testing.T) {
	tests := []struct {
		url     string
		want    string
		wantErr bool
	}{
		{
			url:     "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit?usp=sharing",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
		{
			url:     "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit?usp=sharing",
			want:    "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got, err := ExtractSheetIdFromUrl(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
