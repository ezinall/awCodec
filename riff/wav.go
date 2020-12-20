package riff

import (
	"awCodec/pcm"
	"encoding/binary"
	"fmt"
	"os"
)

var riffId = [4]byte{'R', 'I', 'F', 'F'}
var formatWave = [4]byte{'W', 'A', 'V', 'E'}

var chunkFmt = [4]byte{'f', 'm', 't', ' '}
var chunkData = [4]byte{'d', 'a', 't', 'a'}

const (
	WaveFormatUnknown                int16 = 0x0000             // Microsoft Corporation
	WaveFormatPcm                    int16 = 0x0001             // Microsoft Corporation
	WaveFormatAdpcm                  int16 = 0x0002             // Microsoft Corporation
	WaveFormatIeeeFloat              int16 = 0x0003             // Microsoft Corporation
	WaveFormatVselp                  int16 = 0x0004             // Compaq Computer Corp.
	WaveFormatIbmCvsd                int16 = 0x0005             // IBM Corporation
	WaveFormatAlaw                   int16 = 0x0006             // Microsoft Corporation
	WaveFormatMulaw                  int16 = 0x0007             // Microsoft Corporation
	WaveFormatDts                    int16 = 0x0008             // Microsoft Corporation
	WaveFormatDrm                    int16 = 0x0009             // Microsoft Corporation
	WaveFormatOkiAdpcm               int16 = 0x0010             // OKI
	WaveFormatDviAdpcm               int16 = 0x0011             // Intel Corporation
	WaveFormatImaAdpcm                     = WaveFormatDviAdpcm // Intel Corporation
	WaveFormatMediaspaceAdpcm        int16 = 0x0012             // Videologic
	WaveFormatSierraAdpcm            int16 = 0x0013             // Sierra Semiconductor Corp
	WaveFormatG723Adpcm              int16 = 0x0014             // Antex Electronics Corporation
	WaveFormatDigistd                int16 = 0x0015             // DSP Solutions, Inc.
	WaveFormatDigifix                int16 = 0x0016             // DSP Solutions, Inc.
	WaveFormatDialogicOkiAdpcm       int16 = 0x0017             // Dialogic Corporation
	WaveFormatMediavisionAdpcm       int16 = 0x0018             // Media Vision, Inc.
	WaveFormatCuCodec                int16 = 0x0019             // Hewlett-Packard Company
	WaveFormatYamahaAdpcm            int16 = 0x0020             // Yamaha Corporation of America
	WaveFormatSonarc                 int16 = 0x0021             // Speech Compression
	WaveFormatDspgroupTruespeech     int16 = 0x0022             // DSP Group, Inc
	WaveFormatEchosc1                int16 = 0x0023             // Echo Speech Corporation
	WaveFormatAudiofileAf36          int16 = 0x0024             // Virtual Music, Inc.
	WaveFormatAptx                   int16 = 0x0025             // Audio Processing Technology
	WaveFormatAudiofileAf10          int16 = 0x0026             // Virtual Music, Inc.
	WaveFormatProsody1612            int16 = 0x0027             // Aculab plc
	WaveFormatLrc                    int16 = 0x0028             // Merging Technologies S.A.
	WaveFormatDolbyAc2               int16 = 0x0030             // Dolby Laboratories
	WaveFormatGsm610                 int16 = 0x0031             // Microsoft Corporation
	WaveFormatMsnaudio               int16 = 0x0032             // Microsoft Corporation
	WaveFormatAntexAdpcme            int16 = 0x0033             // Antex Electronics Corporation
	WaveFormatControlResVqlpc        int16 = 0x0034             // Control Resources Limited
	WaveFormatDigireal               int16 = 0x0035             // DSP Solutions, Inc.
	WaveFormatDigiadpcm              int16 = 0x0036             // DSP Solutions, Inc.
	WaveFormatControlResCr10         int16 = 0x0037             // Control Resources Limited
	WaveFormatNmsVbxadpcm            int16 = 0x0038             // Natural MicroSystems
	WaveFormatCsImaadpcm             int16 = 0x0039             // Crystal Semiconductor IMA ADPCM
	WaveFormatEchosc3                int16 = 0x003A             // Echo Speech Corporation
	WaveFormatRockwellAdpcm          int16 = 0x003B             // Rockwell International
	WaveFormatRockwellDigitalk       int16 = 0x003C             // Rockwell International
	WaveFormatXebec                  int16 = 0x003D             // Xebec Multimedia Solutions Limited
	WaveFormatG721Adpcm              int16 = 0x0040             // Antex Electronics Corporation
	WaveFormatG728Celp               int16 = 0x0041             // Antex Electronics Corporation
	WaveFormatMsg723                 int16 = 0x0042             // Microsoft Corporation
	WaveFormatMpeg                   int16 = 0x0050             // Microsoft Corporation
	WaveFormatRt24                   int16 = 0x0052             // InSoft, Inc.
	WaveFormatPac                    int16 = 0x0053             // InSoft, Inc.
	WaveFormatMpeglayer3             int16 = 0x0055             // ISO/MPEG Layer3 Format Tag
	WaveFormatLucentG723             int16 = 0x0059             // Lucent Technologies
	WaveFormatCirrus                 int16 = 0x0060             // Cirrus Logic
	WaveFormatEspcm                  int16 = 0x0061             // ESS Technology
	WaveFormatVoxware                int16 = 0x0062             // Voxware Inc
	WaveFormatCanopusAtrac           int16 = 0x0063             // Canopus, co., Ltd.
	WaveFormatG726Adpcm              int16 = 0x0064             // APICOM
	WaveFormatG722Adpcm              int16 = 0x0065             // APICOM
	WaveFormatDsatDisplay            int16 = 0x0067             // Microsoft Corporation
	WaveFormatVoxwareByteAligned     int16 = 0x0069             // Voxware Inc
	WaveFormatVoxwareAc8             int16 = 0x0070             // Voxware Inc
	WaveFormatVoxwareAc10            int16 = 0x0071             // Voxware Inc
	WaveFormatVoxwareAc16            int16 = 0x0072             // Voxware Inc
	WaveFormatVoxwareAc20            int16 = 0x0073             // Voxware Inc
	WaveFormatVoxwareRt24            int16 = 0x0074             // Voxware Inc
	WaveFormatVoxwareRt29            int16 = 0x0075             // Voxware Inc
	WaveFormatVoxwareRt29hw          int16 = 0x0076             // Voxware Inc
	WaveFormatVoxwareVr12            int16 = 0x0077             // Voxware Inc
	WaveFormatVoxwareVr18            int16 = 0x0078             // Voxware Inc
	WaveFormatVoxwareTq40            int16 = 0x0079             // Voxware Inc
	WaveFormatSoftsound              int16 = 0x0080             // Softsound, Ltd.
	WaveFormatVoxwareTq60            int16 = 0x0081             // Voxware Inc
	WaveFormatMsrt24                 int16 = 0x0082             // Microsoft Corporation
	WaveFormatG729a                  int16 = 0x0083             // AT&amp;T Labs, Inc.
	WaveFormatMviMvi2                int16 = 0x0084             // Motion Pixels
	WaveFormatDfG726                 int16 = 0x0085             // DataFusion Systems (Pty) (Ltd)
	WaveFormatDfGsm610               int16 = 0x0086             // DataFusion Systems (Pty) (Ltd)
	WaveFormatIsiaudio               int16 = 0x0088             // Iterated Systems, Inc.
	WaveFormatOnlive                 int16 = 0x0089             // OnLive! Technologies, Inc.
	WaveFormatSbc24                  int16 = 0x0091             // Siemens Business Communications Sys
	WaveFormatDolbyAc3Spdif          int16 = 0x0092             // Sonic Foundry
	WaveFormatMediasonicG723         int16 = 0x0093             // MediaSonic
	WaveFormatProsody8kbps           int16 = 0x0094             // Aculab plc
	WaveFormatZyxelAdpcm             int16 = 0x0097             // ZyXEL Communications, Inc.
	WaveFormatPhilipsLpcbb           int16 = 0x0098             // Philips Speech Processing
	WaveFormatPacked                 int16 = 0x0099             // Studer Professional Audio AG
	WaveFormatMaldenPhonytalk        int16 = 0x00A0             // Malden Electronics Ltd.
	WaveFormatRhetorexAdpcm          int16 = 0x0100             // Rhetorex Inc.
	WaveFormatIrat                   int16 = 0x0101             // BeCubed Software Inc.
	WaveFormatVivoG723               int16 = 0x0111             // Vivo Software
	WaveFormatVivoSiren              int16 = 0x0112             // Vivo Software
	WaveFormatDigitalG723            int16 = 0x0123             // Digital Equipment Corporation
	WaveFormatSanyoLdAdpcm           int16 = 0x0125             // Sanyo Electric Co., Ltd.
	WaveFormatSiprolabAceplnet       int16 = 0x0130             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabAcelp4800      int16 = 0x0131             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabAcelp8v3       int16 = 0x0132             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabG729           int16 = 0x0133             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabG729a          int16 = 0x0134             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabKelvin         int16 = 0x0135             // Sipro Lab Telecom Inc.
	WaveFormatG726adpcm              int16 = 0x0140             // Dictaphone Corporation
	WaveFormatQualcommPurevoice      int16 = 0x0150             // Qualcomm, Inc.
	WaveFormatQualcommHalfrate       int16 = 0x0151             // Qualcomm, Inc.
	WaveFormatTubgsm                 int16 = 0x0155             // Ring Zero Systems, Inc.
	WaveFormatMsaudio1               int16 = 0x0160             // Microsoft Corporation
	WaveFormatUnisysNapAdpcm         int16 = 0x0170             // Unisys Corp.
	WaveFormatUnisysNapUlaw          int16 = 0x0171             // Unisys Corp.
	WaveFormatUnisysNapAlaw          int16 = 0x0172             // Unisys Corp.
	WaveFormatUnisysNap16k           int16 = 0x0173             // Unisys Corp.
	WaveFormatCreativeAdpcm          int16 = 0x0200             // Creative Labs, Inc
	WaveFormatCreativeFastspeech8    int16 = 0x0202             // Creative Labs, Inc
	WaveFormatCreativeFastspeech10   int16 = 0x0203             // Creative Labs, Inc
	WaveFormatUherAdpcm              int16 = 0x0210             // UHER informatic GmbH
	WaveFormatQuarterdeck            int16 = 0x0220             // Quarterdeck Corporation
	WaveFormatIlinkVc                int16 = 0x0230             // I-link Worldwide
	WaveFormatRawSport               int16 = 0x0240             // Aureal Semiconductor
	WaveFormatEsstAc3                int16 = 0x0241             // ESS Technology, Inc.
	WaveFormatIpiHsx                 int16 = 0x0250             // Interactive Products, Inc.
	WaveFormatIpiRpelp               int16 = 0x0251             // Interactive Products, Inc.
	WaveFormatCs2                    int16 = 0x0260             // Consistent Software
	WaveFormatSonyScx                int16 = 0x0270             // Sony Corp.
	WaveFormatFmTownsSnd             int16 = 0x0300             // Fujitsu Corp.
	WaveFormatBtvDigital             int16 = 0x0400             // Brooktree Corporation
	WaveFormatQdesignMusic           int16 = 0x0450             // QDesign Corporation
	WaveFormatVmeVmpcm               int16 = 0x0680             // AT&amp;T Labs, Inc.
	WaveFormatTpc                    int16 = 0x0681             // AT&amp;T Labs, Inc.
	WaveFormatOligsm                 int16 = 0x1000             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatOliadpcm               int16 = 0x1001             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatOlicelp                int16 = 0x1002             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatOlisbc                 int16 = 0x1003             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatOliopr                 int16 = 0x1004             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatLhCodec                int16 = 0x1100             // Lernout &amp; Hauspie
	WaveFormatNorris                 int16 = 0x1400             // Norris Communications, Inc.
	WaveFormatSoundspaceMusicompress int16 = 0x1500             // AT&amp;T Labs, Inc.
	WaveFormatDvm                    int16 = 0x2000             // FAST Multimedia AG
)

type wave struct {
	ckID   [4]byte
	ckSize int32
	format [4]byte // Format.

	subchunk1Id     [4]byte
	subchunk1Size   int32
	wFormatTag      int16 // Format category.
	nChannels       int16 // Number of channels. 1 for mono or 2 for stereo.
	nSamplesPerSec  int32 // Sampling rate.
	nAvgBytesPerSec int32 // For buffer estimation. sampleRate * nChannels * nBitsPerSample/8.
	nBlockAlign     int16 // Data block size. nChannels * nBitsPerSample/8.

	nBitsPerSample int16 // Sample size.
	subchunk2Id    [4]byte
	subchunk2Size  int32
}

func EncodeWav(samples pcm.Samples, waveFormat int16) {
	context := samples.Context()

	riff := wave{
		ckID:            riffId,
		format:          formatWave,
		subchunk1Id:     chunkFmt,
		subchunk1Size:   16,
		wFormatTag:      waveFormat,
		nChannels:       int16(context.Channels),
		nSamplesPerSec:  int32(context.SampleRate),
		nAvgBytesPerSec: int32(context.SampleRate * context.Channels * samples.BitPerSample() / 8),
		nBlockAlign:     int16(context.Channels * samples.BitPerSample() / 8),
		nBitsPerSample:  int16(samples.BitPerSample()),
		subchunk2Id:     chunkData,
		subchunk2Size:   int32(samples.Len() * samples.BitPerSample() / 8),
	}
	riff.ckSize = riff.subchunk1Size + riff.subchunk2Size

	outFile, _ := os.Create("out.wav")
	_ = binary.Write(outFile, binary.LittleEndian, riff)

	switch v := samples.Pcm().(type) {
	case []int16:
		for _, sample := range v {
			if err := binary.Write(outFile, binary.LittleEndian, sample); err != nil {
				fmt.Println(err)
			}
		}

	case []float32:
		for _, sample := range v {
			if err := binary.Write(outFile, binary.LittleEndian, sample); err != nil {
				fmt.Println(err)
			}
		}

	case []uint8:
		for _, sample := range v {
			if err := binary.Write(outFile, binary.LittleEndian, sample); err != nil {
				fmt.Println(err)
			}
		}
	}
}
