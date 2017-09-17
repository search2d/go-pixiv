package pixiv

import "io/ioutil"

func fixture(path string) []byte {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return buf
}
