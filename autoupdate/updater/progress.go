package updater

import "io"

type progressReader struct {
	io.Reader
	progressCallback func(int64, int64)
	contentLength    int64
	downloaded       int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.downloaded += int64(n)

	if err != io.EOF {
		if pr.contentLength > 0 {
			pr.progressCallback(pr.downloaded, pr.contentLength)
		} else {
			pr.progressCallback(pr.downloaded, 0)
		}
	} else {
		pr.progressCallback(pr.downloaded, pr.downloaded)
	}

	return n, err
}
