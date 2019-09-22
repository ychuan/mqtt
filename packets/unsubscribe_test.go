package packets

import (
	"github.com/stretchr/testify/require"
	"testing"

	"bytes"

	"github.com/jinzhu/copier"
)

func TestUnsubscribeEncode(t *testing.T) {
	require.Contains(t, expectedPackets, Unsubscribe)
	for i, wanted := range expectedPackets[Unsubscribe] {

		if !encodeTestOK(wanted) {
			continue
		}

		require.Equal(t, uint8(10), Unsubscribe, "Incorrect Packet Type [i:%d] %s", i, wanted.desc)
		pk := new(UnsubscribePacket)
		copier.Copy(pk, wanted.packet.(*UnsubscribePacket))

		require.Equal(t, Unsubscribe, pk.Type, "Mismatched Packet Type [i:%d] %s", i, wanted.desc)
		require.Equal(t, Unsubscribe, pk.FixedHeader.Type, "Mismatched FixedHeader Type [i:%d] %s", i, wanted.desc)

		var b bytes.Buffer
		err := pk.Encode(&b)
		if wanted.expect != nil {

			require.Error(t, err, "Expected error writing buffer [i:%d] %s", i, wanted.desc)

		} else {

			require.NoError(t, err, "Error writing buffer [i:%d] %s", i, wanted.desc)

			encoded := b.Bytes()

			require.Equal(t, len(wanted.rawBytes), len(encoded), "Mismatched packet length [i:%d] %s", i, wanted.desc)
			if wanted.meta != nil {
				require.Equal(t, byte(Unsubscribe<<4)|wanted.meta.(byte), encoded[0], "Mismatched fixed header bytes [i:%d] %s", i, wanted.desc)
			} else {
				require.Equal(t, byte(Unsubscribe<<4), encoded[0], "Mismatched fixed header bytes [i:%d] %s", i, wanted.desc)
			}

			require.NoError(t, err, "Error writing buffer [i:%d] %s", i, wanted.desc)
			require.EqualValues(t, wanted.rawBytes, encoded, "Mismatched byte values [i:%d] %s", i, wanted.desc)

			require.Equal(t, wanted.packet.(*UnsubscribePacket).PacketID, pk.PacketID, "Mismatched Packet ID [i:%d] %s", i, wanted.desc)
			require.Equal(t, wanted.packet.(*UnsubscribePacket).Topics, pk.Topics, "Mismatched Topics slice [i:%d] %s", i, wanted.desc)

		}

	}

}

func TestUnsubscribeDecode(t *testing.T) {

	require.Contains(t, expectedPackets, Unsubscribe)
	for i, wanted := range expectedPackets[Unsubscribe] {

		if !decodeTestOK(wanted) {
			continue
		}

		require.Equal(t, uint8(10), Unsubscribe, "Incorrect Packet Type [i:%d] %s", i, wanted.desc)

		pk := newPacket(Unsubscribe).(*UnsubscribePacket)
		err := pk.Decode(wanted.rawBytes[2:]) // Unpack skips fixedheader.
		if wanted.failFirst != nil {
			require.Error(t, err, "Expected error unpacking buffer [i:%d] %s", i, wanted.desc)
			require.Equal(t, wanted.failFirst, err.Error(), "Expected fail state; %v [i:%d] %s", err.Error(), i, wanted.desc)
			continue
		}

		require.NoError(t, err, "Error unpacking buffer [i:%d] %s", i, wanted.desc)

		require.Equal(t, wanted.packet.(*UnsubscribePacket).PacketID, pk.PacketID, "Mismatched Packet ID [i:%d] %s", i, wanted.desc)
		require.Equal(t, wanted.packet.(*UnsubscribePacket).Topics, pk.Topics, "Mismatched Topics slice [i:%d] %s", i, wanted.desc)
	}

}

func BenchmarkUnsubscribeDecode(b *testing.B) {
	pk := newPacket(Unsubscribe).(*UnsubscribePacket)
	pk.FixedHeader.decode(expectedPackets[Unsubscribe][0].rawBytes[0])

	for n := 0; n < b.N; n++ {
		pk.Decode(expectedPackets[Unsubscribe][0].rawBytes[2:])
	}
}
