package iccprof

// Copied from https://github.com/vimeo/go-iccjpeg

import (
	"bufio"
	"errors"
	"io"
)

const (
	/* JPEG Markers. */
	soiMarker  = 0xD8
	eoiMarker  = 0xD9
	app2Marker = 0xE2
	rst0Marker = 0xD0
	rst7Marker = 0xD7

	/* Others. */
	iccHeaderLen = 14
)

func getSize(input io.Reader) (int, error) {
	var buf [2]byte
	_, err := io.ReadFull(input, buf[0:2])
	if err != nil {
		return 0, err
	}

	ret := int(buf[0])<<8 + int(buf[1]) - 2
	if ret < 0 {
		return ret, errors.New("invalid segment length")
	}

	return ret, nil
}

func GetICCBuf(input io.Reader) ([]byte, error) {
	var buf [1024]byte
	var err error
	in := bufio.NewReader(input)

	_, err = io.ReadFull(in, buf[0:2])
	if err != nil {
		return nil, err
	} else if buf[0] != 0xFF && buf[1] != soiMarker {
		return nil, errors.New("no SOI Marker")
	}

	var icc_data [][]byte
	icc_length := 0
	read_profs := 0
	num_markers := -1
	for {
		_, err = io.ReadFull(in, buf[0:2])
		if err != nil {
			return nil, err
		}

		/* Handle broken jpegs. */
		for buf[0] != 0xFF {
			buf[0] = buf[1]
			buf[1], err = in.ReadByte()
			if err != nil {
				return nil, err
			}
		}

		/* Skip 00 markers. */
		if buf[1] == 0 {
			continue
		}

		/* Skip stuffing. */
		for buf[1] == 0xFF {
			buf[1], err = in.ReadByte()
			if err != nil {
				return nil, err
			}
		}

		/* We reached the end of the image. */
		if buf[1] == eoiMarker {
			break
		}

		/* Are we at an APP2 marker? */
		if buf[1] != app2Marker {
			/* Skip RST if need be. */
			if buf[1] >= rst0Marker && buf[1] <= rst7Marker {
				continue
			}

			size, err := getSize(in)
			if err != nil {
				return nil, err
			}

			/* Skip non-APP2. */
			_, err = io.CopyN(io.Discard, in, int64(size))
			if err != nil {
				return nil, err
			}
			continue
		}

		size, err := getSize(in)
		if err != nil {
			return nil, err
		} else if size < iccHeaderLen {
			return nil, errors.New("ICC segment invalid")
		}

		_, err = io.ReadFull(in, buf[0:12])
		if err != nil {
			return nil, err
		}

		if string(buf[0:11]) != "ICC_PROFILE" || buf[11] != 0 {
			return nil, errors.New("ICC segment invalid")
		}

		seqno, err := in.ReadByte()
		if err != nil {
			return nil, err
		} else if seqno == 0 {
			return nil, errors.New("invalid sequence number")
		}

		num, err := in.ReadByte()
		if err != nil {
			return nil, err
		} else if num_markers == -1 {
			num_markers = int(num)
			icc_data = make([][]byte, num_markers)
		} else if int(num) != num_markers {
			return nil, errors.New("invalid ICC segment (num_markers != cur_num_markers)")
		}

		if int(seqno) > num_markers {
			return nil, errors.New("invalid ICC segment (seqno > num_markers)")
		}

		icc_data[seqno-1] = make([]byte, size-iccHeaderLen)
		_, err = io.ReadFull(in, icc_data[seqno-1])
		if err != nil {
			return nil, err
		}

		icc_length += size - iccHeaderLen
		read_profs++

		if read_profs == num_markers {
			ret := make([]byte, 0, icc_length)
			for _, data := range icc_data {
				ret = append(ret, data...)
			}
			return ret, nil
		}
	}

	return nil, nil
}
