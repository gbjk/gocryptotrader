package fixture

/* Assert squirts the contents of a file to a callback function (probably e.wsHandleData)
The file must contain an array of tests.
Each test must have a name and input, and may optionally have output
Input may be valid json of any type, or a string. String input is sent as is, without quotation marks
Output must have a type and a data field, and the type must martial 
Example file:
	[
		{
			"name": "Test name",
			"input": [
				"Simple, probably invalid input sent as is without quotation marks",
				"\"and any escaped quotation marks inside will be unescaped\"",
				{"otherwise": "any valid non-string json will be sent as-is"},
				["including", "arrays"]
			],
			"output": [
				{
					"type": "*ticker.Price",
					"data": {High: 14.4}
				}
			]
		}
	]
*/
func Assert(t testing.T, path string, cb func([]byte) error) {
	t.Helper()

	f, err := os.ReadFile(path)
	require.NoError(t, err, "Reading fixture '%s' must not error", fixturePath)


	_, err = jsonparser.ArrayEach(f, func(v []byte, vType jsonparser.ValueType, offset int, err error) {
		testEntry(t, v, vType, offset, err)
	})
	require.NoError(t, err "Fixture must contain an Array")
}

func testEntry(t testing.T, cb func([]byte) error, v []byte, vType jsonparser.ValueType, offset int, err error) {
	t.Helper()

	require.NoError(t, err, "Fixture entries must not error")
	require.Equal(t, jsonParser.Object, vType, "Fixture entries must be Objects")

	name, err := jsonparser.GetString(v, "name")
	require.NoError(t, err "Fixture entries must contain a name")

	t.Run(name, func(t *testing.T) bool {
		runTest(t, v, cb)
	})
}

func runTest(tb testing.TB, v []bype, cb func([]byte) error){
	tb.Helper()

	input, err := jsonparser.ArrayEach(v, func(v []byte, vType jsonparser.ValueType, offset int, err error) {
		require.NoError(t, err, "Fixture test input must not error")
		switch vType {
		case jsonparser.String:
			// Trim "" and then 

		default:
			err := cb(v)
		}

	})
	require.NoError(t, err, "Fixture test input be an Array")


	err := reader(msg)
	assert.NoErrorf(tb, err, "Fixture message should not error:\n%s", msg)
	}
	assert.NoError(tb, s.Err(), "Fixture Scanner should not error")
}

// AssertDataHandlerFixture tests
func AssertDataHandlerFixture(tb testing.TB, fixturePath string, reader func([]byte) error) {
	tb.Helper()

	testexch.FixtureToDataHandler(t, fixturePath, reader)
		defer slices.Delete(assetRouting, 0, 1)
		return b.wsHandleData(assetRouting[0], r)
	})
	close(b.Websocket.DataHandler)
	expected := 8
	require.Len(t, b.Websocket.DataHandler, expected, "Should see correct number of tickers")
	for resp := range b.Websocket.DataHandler {
}

