package headers

type Headers map[string]string

func NewHeaders() Headers {
	headers := make(Headers)
	return headers
}

func (h Headers) parse(data []byte) (int, bool, error) {
read := 0

	for {
		switch r.state {
		case StateInit:
            rl, n, err := parseRequestLine(msg[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				return read, nil
			}

			r.RequestLine = *rl
			read += n

			r.state = StateDone

		case StateDone:
			return read, nil
		case StateError:
			return 0, ErrorRequestInErrorState
		}
	}
	return 0, false, nil
}
