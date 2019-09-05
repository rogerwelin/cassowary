package main

const (
	summaryTable = `` + "\n\n" +
		` DNS Lookup......................: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		` TCP Connect.....................: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		` Server Processing...............: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		` Content Transfer................: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		`` + "\n" +
		`Summary: ` + "\n" +
		` Total Req.......................: %s` + "\n" +
		` Failed Req......................: %s` + "\n" +
		` Req/s...........................: %s` + "\n\n"

	summaryTLSTable = `` + "\n\n" +
		` DNS Lookup......................: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		` TCP Connect.....................: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		` TLS Handshake...................: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		` Server Processing...............: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		` Content Transfer................: Avg/mean=%sms ` + "\t" + `Median=%sms ` + "\t" + `p(95)=%sms` + "\n" +
		`` + "\n" +
		`Summary: ` + "\n" +
		` Total Req.......................: %s` + "\n" +
		` Failed Req......................: %s` + "\n" +
		` Req/s...........................: %s` + "\n\n"
)

func (c *cassowary) coordinate() error {

	if c.fileMode {
		urlSuffixes, err := readFile(c.inputFile)
		if err != nil {
			return err
		}
	}

	return nil
}
