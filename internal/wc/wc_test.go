package wc

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessor_reset(t *testing.T) {
	t.Parallel()

	p := &Processor{
		lines: 1,
		words: 2,
		chars: 3,
		bytes: 4,
	}
	p.reset()

	require.Equal(t, 0, p.Lines())
	require.Equal(t, 0, p.Words())
	require.Equal(t, 0, p.Chars())
	require.Equal(t, 0, p.Bytes())
}

func TestProcessor_Analyze(t *testing.T) {
	t.Parallel()

	type want struct {
		bytes int
		chars int
		lines int
		words int
	}

	tests := []struct {
		name    string
		proc    *Processor
		fd      *os.File
		want    want
		wantErr error
	}{
		{
			name:    "failure - input is nil",
			wantErr: fmt.Errorf("input cannot be nil"),
		},
		{
			name: "failure - directory",
			proc: &Processor{},
			fd: func() *os.File {
				f, err := os.Open("testdata")
				require.NoError(t, err)
				return f
			}(),
			wantErr: fmt.Errorf("read testdata: is a directory"),
		},
		{
			name: "success - bytes only (file)",
			proc: &Processor{cfg: &Config{CountBytes: true}},
			fd: func() *os.File {
				f, err := os.Open("testdata/test.txt")
				require.NoError(t, err)
				return f
			}(),
			want: want{bytes: 342190},
		},
		{
			name: "success - lines, words, bytes, chars (file)",
			proc: &Processor{
				cfg: &Config{
					CountLines: true,
					CountWords: true,
					CountBytes: true,
					CountChars: true,
				},
			},
			fd: func() *os.File {
				f, err := os.Open("testdata/test.txt")
				require.NoError(t, err)
				return f
			}(),
			want: want{
				bytes: 342190,
				words: 58164,
				lines: 7145,
				chars: 339292,
			},
		},
	}

	for _, tt := range tests {
		tt := tt // retain for parallel

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.proc.Analyze(tt.fd)
			if tt.fd != nil {
				tt.fd.Close()
			}

			if tt.wantErr != nil {
				require.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.bytes, tt.proc.Bytes())
				require.Equal(t, tt.want.chars, tt.proc.Chars())
				require.Equal(t, tt.want.lines, tt.proc.Lines())
				require.Equal(t, tt.want.words, tt.proc.Words())
			}
		})
	}
}
