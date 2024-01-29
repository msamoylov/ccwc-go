package file

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    int64
		wantErr error
	}{
		{
			name:    "failure - blank path",
			wantErr: errors.New("path cannot be blank"),
		},
		{
			name:    "failure - path does not exist",
			path:    "test",
			wantErr: errors.New("open test: no such file or directory"),
		},
		{
			name:    "failure - path is a directory",
			path:    ".",
			wantErr: errors.New("read .: is a directory"),
		},
		{
			name: "success - file exists",
			path: "testdata/test.txt",
			want: 342190,
		},
	}

	for _, tt := range tests {
		tt := tt // retain for parallel

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Size(tt.path)
			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
