package riff

import (
	"awCodec/pcm"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	WaveFormatUnknown                 uint16 = 0x0000             // Microsoft Corporation
	WaveFormatPcm                     uint16 = 0x0001             // Microsoft Corporation
	WaveFormatAdpcm                   uint16 = 0x0002             // Microsoft Corporation
	WaveFormatIeeeFloat               uint16 = 0x0003             // Microsoft Corporation
	WaveFormatVselp                   uint16 = 0x0004             // Compaq Computer Corp.
	WaveFormatIbmCvsd                 uint16 = 0x0005             // IBM Corporation
	WaveFormatAlaw                    uint16 = 0x0006             // Microsoft Corporation
	WaveFormatMulaw                   uint16 = 0x0007             // Microsoft Corporation
	WaveFormatDts                     uint16 = 0x0008             // Microsoft Corporation
	WaveFormatDrm                     uint16 = 0x0009             // Microsoft Corporation
	WaveFormatWmavoice9               uint16 = 0x000A             // Microsoft Corporation
	WaveFormatWmavoice10              uint16 = 0x000B             // Microsoft Corporation
	WaveFormatOkiAdpcm                uint16 = 0x0010             // OKI
	WaveFormatDviAdpcm                uint16 = 0x0011             // Intel Corporation
	WaveFormatImaAdpcm                       = WaveFormatDviAdpcm // Intel Corporation
	WaveFormatMediaspaceAdpcm         uint16 = 0x0012             // Videologic
	WaveFormatSierraAdpcm             uint16 = 0x0013             // Sierra Semiconductor Corp
	WaveFormatG723Adpcm               uint16 = 0x0014             // Antex Electronics Corporation
	WaveFormatDigistd                 uint16 = 0x0015             // DSP Solutions, Inc.
	WaveFormatDigifix                 uint16 = 0x0016             // DSP Solutions, Inc.
	WaveFormatDialogicOkiAdpcm        uint16 = 0x0017             // Dialogic Corporation
	WaveFormatMediavisionAdpcm        uint16 = 0x0018             // Media Vision, Inc.
	WaveFormatCuCodec                 uint16 = 0x0019             // Hewlett-Packard Company
	WaveFormatHpDynVoice              uint16 = 0x001A             // Hewlett-Packard Company
	WaveFormatYamahaAdpcm             uint16 = 0x0020             // Yamaha Corporation of America
	WaveFormatSonarc                  uint16 = 0x0021             // Speech Compression
	WaveFormatDspgroupTruespeech      uint16 = 0x0022             // DSP Group, Inc
	WaveFormatEchosc1                 uint16 = 0x0023             // Echo Speech Corporation
	WaveFormatAudiofileAf36           uint16 = 0x0024             // Virtual Music, Inc.
	WaveFormatAptx                    uint16 = 0x0025             // Audio Processing Technology
	WaveFormatAudiofileAf10           uint16 = 0x0026             // Virtual Music, Inc.
	WaveFormatProsody1612             uint16 = 0x0027             // Aculab plc
	WaveFormatLrc                     uint16 = 0x0028             // Merging Technologies S.A.
	WaveFormatDolbyAc2                uint16 = 0x0030             // Dolby Laboratories
	WaveFormatGsm610                  uint16 = 0x0031             // Microsoft Corporation
	WaveFormatMsnaudio                uint16 = 0x0032             // Microsoft Corporation
	WaveFormatAntexAdpcme             uint16 = 0x0033             // Antex Electronics Corporation
	WaveFormatControlResVqlpc         uint16 = 0x0034             // Control Resources Limited
	WaveFormatDigireal                uint16 = 0x0035             // DSP Solutions, Inc.
	WaveFormatDigiadpcm               uint16 = 0x0036             // DSP Solutions, Inc.
	WaveFormatControlResCr10          uint16 = 0x0037             // Control Resources Limited
	WaveFormatNmsVbxadpcm             uint16 = 0x0038             // Natural MicroSystems
	WaveFormatCsImaadpcm              uint16 = 0x0039             // Crystal Semiconductor IMA ADPCM
	WaveFormatEchosc3                 uint16 = 0x003A             // Echo Speech Corporation
	WaveFormatRockwellAdpcm           uint16 = 0x003B             // Rockwell International
	WaveFormatRockwellDigitalk        uint16 = 0x003C             // Rockwell International
	WaveFormatXebec                   uint16 = 0x003D             // Xebec Multimedia Solutions Limited
	WaveFormatG721Adpcm               uint16 = 0x0040             // Antex Electronics Corporation
	WaveFormatG728Celp                uint16 = 0x0041             // Antex Electronics Corporation
	WaveFormatMsg723                  uint16 = 0x0042             // Microsoft Corporation
	WaveFormatIntelG7231              uint16 = 0x0043             // Intel Corp.
	WaveFormatIntelG729               uint16 = 0x0044             // Intel Corp.
	WaveFormatSharpG726               uint16 = 0x0045             // Sharp
	WaveFormatMpeg                    uint16 = 0x0050             // Microsoft Corporation
	WaveFormatRt24                    uint16 = 0x0052             // InSoft, Inc.
	WaveFormatPac                     uint16 = 0x0053             // InSoft, Inc.
	WaveFormatMpegLayer3              uint16 = 0x0055             // ISO/MPEG Layer3 Format Tag
	WaveFormatLucentG723              uint16 = 0x0059             // Lucent Technologies
	WaveFormatCirrus                  uint16 = 0x0060             // Cirrus Logic
	WaveFormatEspcm                   uint16 = 0x0061             // ESS Technology
	WaveFormatVoxware                 uint16 = 0x0062             // Voxware Inc
	WaveFormatCanopusAtrac            uint16 = 0x0063             // Canopus, co., Ltd.
	WaveFormatG726Adpcm               uint16 = 0x0064             // APICOM
	WaveFormatG722Adpcm               uint16 = 0x0065             // APICOM
	WaveFormatDsatDisplay             uint16 = 0x0067             // Microsoft Corporation
	WaveFormatVoxwareByteAligned      uint16 = 0x0069             // Voxware Inc
	WaveFormatVoxwareAc8              uint16 = 0x0070             // Voxware Inc
	WaveFormatVoxwareAc10             uint16 = 0x0071             // Voxware Inc
	WaveFormatVoxwareAc16             uint16 = 0x0072             // Voxware Inc
	WaveFormatVoxwareAc20             uint16 = 0x0073             // Voxware Inc
	WaveFormatVoxwareRt24             uint16 = 0x0074             // Voxware Inc
	WaveFormatVoxwareRt29             uint16 = 0x0075             // Voxware Inc
	WaveFormatVoxwareRt29hw           uint16 = 0x0076             // Voxware Inc
	WaveFormatVoxwareVr12             uint16 = 0x0077             // Voxware Inc
	WaveFormatVoxwareVr18             uint16 = 0x0078             // Voxware Inc
	WaveFormatVoxwareTq40             uint16 = 0x0079             // Voxware Inc
	WaveFormatVoxwareSc3              uint16 = 0x007A             // Voxware Inc
	WaveFormatVoxwareSc31             uint16 = 0x007B             // Voxware Inc
	WaveFormatSoftsound               uint16 = 0x0080             // Softsound, Ltd.
	WaveFormatVoxwareTq60             uint16 = 0x0081             // Voxware Inc
	WaveFormatMsrt24                  uint16 = 0x0082             // Microsoft Corporation
	WaveFormatG729a                   uint16 = 0x0083             // AT&amp;T Labs, Inc.
	WaveFormatMviMvi2                 uint16 = 0x0084             // Motion Pixels
	WaveFormatDfG726                  uint16 = 0x0085             // DataFusion Systems (Pty) (Ltd)
	WaveFormatDfGsm610                uint16 = 0x0086             // DataFusion Systems (Pty) (Ltd)
	WaveFormatIsiaudio                uint16 = 0x0088             // Iterated Systems, Inc.
	WaveFormatOnlive                  uint16 = 0x0089             // OnLive! Technologies, Inc.
	WaveFormatMultitudeFtSx20         uint16 = 0x008A             // Multitude Inc.
	WaveFormatInfocomItsG721Adpcm     uint16 = 0x008B             // Infocom
	WaveFormatConvediaG729            uint16 = 0x008C             // Convedia Corp.
	WaveFormatCongruency              uint16 = 0x008D             // Congruency Inc.
	WaveFormatSbc24                   uint16 = 0x0091             // Siemens Business Communications Sys
	WaveFormatDolbyAc3Spdif           uint16 = 0x0092             // Sonic Foundry
	WaveFormatMediasonicG723          uint16 = 0x0093             // MediaSonic
	WaveFormatProsody8kbps            uint16 = 0x0094             // Aculab plc
	WaveFormatZyxelAdpcm              uint16 = 0x0097             // ZyXEL Communications, Inc.
	WaveFormatPhilipsLpcbb            uint16 = 0x0098             // Philips Speech Processing
	WaveFormatPacked                  uint16 = 0x0099             // Studer Professional Audio AG
	WaveFormatMaldenPhonytalk         uint16 = 0x00A0             // Malden Electronics Ltd.
	WaveFormatRacalRecorderGsm        uint16 = 0x00A1             // Racal recorders
	WaveFormatRacalRecorderG720A      uint16 = 0x00A2             // Racal recorders
	WaveFormatRacalRecorderG7231      uint16 = 0x00A3             // Racal recorders
	WaveFormatRacalRecorderTetraAcelp uint16 = 0x00A4             // Racal recorders
	WaveFormatNecAac                  uint16 = 0x00B0             // NEC Corp.
	WaveFormatRawAac1                 uint16 = 0x00FF             // For Raw AAC, with format block AudioSpecificConfig() (as defined by MPEG-4), that follows WAVEFORMATEX
	WaveFormatRhetorexAdpcm           uint16 = 0x0100             // Rhetorex Inc.
	WaveFormatIrat                    uint16 = 0x0101             // BeCubed Software Inc.
	WaveFormatVivoG723                uint16 = 0x0111             // Vivo Software
	WaveFormatVivoSiren               uint16 = 0x0112             // Vivo Software
	WaveFormatPhilipsCelp             uint16 = 0x0120             // Philips Speech Processing
	WaveFormatPhilipsGrundig          uint16 = 0x0121             // Philips Speech Processing
	WaveFormatDigitalG723             uint16 = 0x0123             // Digital Equipment Corporation
	WaveFormatSanyoLdAdpcm            uint16 = 0x0125             // Sanyo Electric Co., Ltd.
	WaveFormatSiprolabAceplnet        uint16 = 0x0130             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabAcelp4800       uint16 = 0x0131             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabAcelp8v3        uint16 = 0x0132             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabG729            uint16 = 0x0133             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabG729a           uint16 = 0x0134             // Sipro Lab Telecom Inc.
	WaveFormatSiprolabKelvin          uint16 = 0x0135             // Sipro Lab Telecom Inc.
	WaveFormatVoiceageAmr             uint16 = 0x0136             // VoiceAge Corp.
	WaveFormatG726adpcm               uint16 = 0x0140             // Dictaphone Corporation
	WaveFormatDictaphoneCelp68        uint16 = 0x0141             // Dictaphone Corporation
	WaveFormatDictaphoneCelp54        uint16 = 0x0142             // Dictaphone Corporation
	WaveFormatQualcommPurevoice       uint16 = 0x0150             // Qualcomm, Inc.
	WaveFormatQualcommHalfrate        uint16 = 0x0151             // Qualcomm, Inc.
	WaveFormatTubgsm                  uint16 = 0x0155             // Ring Zero Systems, Inc.
	WaveFormatMsaudio1                uint16 = 0x0160             // Microsoft Corporation
	WaveFormatWmaudio2                uint16 = 0x0161             // Microsoft Corporation
	WaveFormatWmaudio3                uint16 = 0x0162             // Microsoft Corporation
	WaveFormatWmaudioLossless         uint16 = 0x0163             // Microsoft Corporation
	WaveFormatWmaspdif                uint16 = 0x0164             // Microsoft Corporation
	WaveFormatUnisysNapAdpcm          uint16 = 0x0170             // Unisys Corp.
	WaveFormatUnisysNapUlaw           uint16 = 0x0171             // Unisys Corp.
	WaveFormatUnisysNapAlaw           uint16 = 0x0172             // Unisys Corp.
	WaveFormatUnisysNap16k            uint16 = 0x0173             // Unisys Corp.
	WaveFormatSycomAcmSyc008          uint16 = 0x0174             // SyCom Technologies
	WaveFormatSycomAcmSyc701G726l     uint16 = 0x0175             // SyCom Technologies
	WaveFormatSycomAcmSyc701Celp54    uint16 = 0x0176             // SyCom Technologies
	WaveFormatSycomAcmSyc701Celp68    uint16 = 0x0177             // SyCom Technologies
	WaveFormatKnowledgeAdventureAdpcm uint16 = 0x0178             // Knowledge Adventure, Inc.
	WaveFormatFraunhoferIisMpeg2Aac   uint16 = 0x0180             // Fraunhofer IIS
	WaveFormatDtsDs                   uint16 = 0x0190             // Digital Theatre Systems, Inc.
	WaveFormatCreativeAdpcm           uint16 = 0x0200             // Creative Labs, Inc
	WaveFormatCreativeFastspeech8     uint16 = 0x0202             // Creative Labs, Inc
	WaveFormatCreativeFastspeech10    uint16 = 0x0203             // Creative Labs, Inc
	WaveFormatUherAdpcm               uint16 = 0x0210             // UHER informatic GmbH
	WaveFormatUleadDvAudio            uint16 = 0x0215             // Ulead Systems, Inc.
	WaveFormatUleadDvAudio1           uint16 = 0x0216             // Ulead Systems, Inc.
	WaveFormatQuarterdeck             uint16 = 0x0220             // Quarterdeck Corporation
	WaveFormatIlinkVc                 uint16 = 0x0230             // I-link Worldwide
	WaveFormatRawSport                uint16 = 0x0240             // Aureal Semiconductor
	WaveFormatEsstAc3                 uint16 = 0x0241             // ESS Technology, Inc.
	WaveFormatGenericPassthru         uint16 = 0x0249             //
	WaveFormatIpiHsx                  uint16 = 0x0250             // Interactive Products, Inc.
	WaveFormatIpiRpelp                uint16 = 0x0251             // Interactive Products, Inc.
	WaveFormatCs2                     uint16 = 0x0260             // Consistent Software
	WaveFormatSonyScx                 uint16 = 0x0270             // Sony Corp.
	WaveFormatSonyScy                 uint16 = 0x0271             // Sony Corp.
	WaveFormatSonyAtrac3              uint16 = 0x0272             // Sony Corp.
	WaveFormatSonySpc                 uint16 = 0x0273             // Sony Corp.
	WaveFormatTelumAudio              uint16 = 0x0280             // Telum Inc.
	WaveFormatTelumIaAudio            uint16 = 0x0281             // Telum Inc.
	WaveFormatNorcomVoiceSystemsAdpcm uint16 = 0x0285             // Norcom Electronics Corp.
	WaveFormatFmTownsSnd              uint16 = 0x0300             // Fujitsu Corp.
	WaveFormatMicronas                uint16 = 0x0350             // Micronas Semiconductors, Inc.
	WaveFormatMicronasCelp833         uint16 = 0x0351             // Micronas Semiconductors, Inc.
	WaveFormatBtvDigital              uint16 = 0x0400             // Brooktree Corporation
	WaveFormatIntelMusicCoder         uint16 = 0x0401             // Intel Corp.
	WaveFormatIndeoAudio              uint16 = 0x0402             // Ligos
	WaveFormatQdesignMusic            uint16 = 0x0450             // QDesign Corporation
	WaveFormatOn2Vp7Audio             uint16 = 0x0500             // On2 Technologies
	WaveFormatOn2Vp6Audio             uint16 = 0x0501             // On2 Technologies
	WaveFormatVmeVmpcm                uint16 = 0x0680             // AT&amp;T Labs, Inc.
	WaveFormatTpc                     uint16 = 0x0681             // AT&amp;T Labs, Inc.
	WaveFormatLightwaveLossless       uint16 = 0x08AE             // Clearjump
	WaveFormatOligsm                  uint16 = 0x1000             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatOliadpcm                uint16 = 0x1001             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatOlicelp                 uint16 = 0x1002             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatOlisbc                  uint16 = 0x1003             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatOliopr                  uint16 = 0x1004             // Ing C. Olivetti &amp; C., S.p.A.
	WaveFormatLhCodec                 uint16 = 0x1100             // Lernout &amp; Hauspie
	WaveFormatLhCodecCelp             uint16 = 0x1101             // Lernout & Hauspie
	WaveFormatLhCodecSbc8             uint16 = 0x1102             // Lernout & Hauspie
	WaveFormatLhCodecSbc12            uint16 = 0x1103             // Lernout & Hauspie
	WaveFormatLhCodecSbc16            uint16 = 0x1104             // Lernout & Hauspie
	WaveFormatNorris                  uint16 = 0x1400             // Norris Communications, Inc.
	WaveFormatIsiaudio2               uint16 = 0x1401             // ISIAudio
	WaveFormatSoundspaceMusicompress  uint16 = 0x1500             // AT&T Labs, Inc.
	WaveFormatMpegAdtsAac             uint16 = 0x1600             // Microsoft Corporation
	WaveFormatMpegRawAac              uint16 = 0x1601             // Microsoft Corporation
	WaveFormatMpegLoas                uint16 = 0x1602             // Microsoft Corporation (MPEG-4 Audio Transport Streams (LOAS/LATM)
	WaveFormatNokiaMpegAdtsAac        uint16 = 0x1608             // Microsoft Corporation
	WaveFormatNokiaMpegRawAac         uint16 = 0x1609             // Microsoft Corporation
	WaveFormatVodafoneMpegAdtsAac     uint16 = 0x160A             // Microsoft Corporation
	WaveFormatVodafoneMpegRawAac      uint16 = 0x160B             // Microsoft Corporation
	WaveFormatMpegHeaac               uint16 = 0x1610             // Microsoft Corporation (MPEG-2 AAC or MPEG-4 HE-AAC v1/v2 streams with any payload (ADTS, ADIF, LOAS/LATM, RAW).
	WaveFormatVoxwareRt24Speech       uint16 = 0x181C             // Voxware Inc.
	WaveFormatSonicfoundryLossless    uint16 = 0x1971             // Sonic Foundry
	WaveFormatInningsTelecomAdpcm     uint16 = 0x1979             // Innings Telecom Inc.
	WaveFormatLucentSx8300p           uint16 = 0x1C07             // Lucent Technologies
	WaveFormatLucentSx5363s           uint16 = 0x1C0C             // Lucent Technologies
	WaveFormatCuseeme                 uint16 = 0x1F03             // CUSeeMe
	WaveFormatNtcsoftAlf2cmAcm        uint16 = 0x1FC4             // NTCSoft
	WaveFormatDvm                     uint16 = 0x2000             // FAST Multimedia AG
	WaveFormatDts2                    uint16 = 0x2001             //
	WaveFormatMakeavis                uint16 = 0x3313             //
	WaveFormatDivioMpeg4Aac           uint16 = 0x4143             // Divio, Inc.
	WaveFormatNokiaAdaptiveMultirate  uint16 = 0x4201             // Nokia
	WaveFormatDivioG726               uint16 = 0x4243             // Divio, Inc.
	WaveFormatLeadSpeech              uint16 = 0x434C             // LEAD Technologies
	WaveFormatLeadVorbis              uint16 = 0x564C             // LEAD Technologies
	WaveFormatWavpackAudio            uint16 = 0x5756             // xiph.org
	WaveFormatAlac                    uint16 = 0x6C61             // Apple Lossless
	WaveFormatOggVorbisMode1          uint16 = 0x674F             // Ogg Vorbis
	WaveFormatOggVorbisMode2          uint16 = 0x6750             // Ogg Vorbis
	WaveFormatOggVorbisMode3          uint16 = 0x6751             // Ogg Vorbis
	WaveFormatOggVorbisMode1Plus      uint16 = 0x676F             // Ogg Vorbis
	WaveFormatOggVorbisMode2Plus      uint16 = 0x6770             // Ogg Vorbis
	WaveFormatOggVorbisMode3Plus      uint16 = 0x6771             // Ogg Vorbis
	WaveFormat3comNbx                 uint16 = 0x7000             // 3COM Corp.
	WaveFormatOpus                    uint16 = 0x704F             // Opus
	WaveFormatFaadAac                 uint16 = 0x706D             //
	WaveFormatAmrNb                   uint16 = 0x7361             // AMR Narrowband
	WaveFormatAmrWb                   uint16 = 0x7362             // AMR Wideband
	WaveFormatAmrWp                   uint16 = 0x7363             // AMR Wideband Plus
	WaveFormatGsmAmrCbr               uint16 = 0x7A21             // GSMA/3GPP
	WaveFormatGsmAmrVbrSid            uint16 = 0x7A22             // GSMA/3GPP
	WaveFormatComverseInfosysG7231    uint16 = 0xA100             // Comverse Infosys
	WaveFormatComverseInfosysAvqsbc   uint16 = 0xA101             // Comverse Infosys
	WaveFormatComverseInfosysSbc      uint16 = 0xA102             // Comverse Infosys
	WaveFormatSymbolG729A             uint16 = 0xA103             // Symbol Technologies
	WaveFormatVoiceageAmrWb           uint16 = 0xA104             // VoiceAge Corp.
	WaveFormatIngenientG726           uint16 = 0xA105             // Ingenient Technologies, Inc.
	WaveFormatMpeg4Aac                uint16 = 0xA106             // ISO/MPEG­4
	WaveFormatEncoreG726              uint16 = 0xA107             // Encore Software
	WaveFormatZollAsao                uint16 = 0xA108             // ZOLL Medical Corp.
	WaveFormatSpeexVoice              uint16 = 0xA109             // xiph.org
	WaveFormatVianixMasc              uint16 = 0xA10A             // Vianix LLC
	WaveFormatWm9SpectrumAnalyzer     uint16 = 0xA10B             // Microsoft
	WaveFormatWmfSpectrumAnayzer      uint16 = 0xA10C             // Microsoft
	WaveFormatGsm610_                 uint16 = 0xA10D             //
	WaveFormatGsm620                  uint16 = 0xA10E             //
	WaveFormatGsm660                  uint16 = 0xA10F             //
	WaveFormatGsm690                  uint16 = 0xA110             //
	WaveFormatGsmAdaptiveMultirateWb  uint16 = 0xA111             //
	WaveFormatPolycomG722             uint16 = 0xA112             // Polycom
	WaveFormatPolycomG728             uint16 = 0xA113             // Polycom
	WaveFormatPolycomG729A            uint16 = 0xA114             // Polycom
	WaveFormatPolycomSiren            uint16 = 0xA115             // Polycom
	WaveFormatGlobalIpIlbc            uint16 = 0xA116             // Global IP
	WaveFormatRadiotimeTimeShiftRadio uint16 = 0xA117             // RadioTime
	WaveFormatNiceAca                 uint16 = 0xA118             // Nice Systems
	WaveFormatNiceAdpcm               uint16 = 0xA119             // Nice Systems
	WaveFormatVocordG721              uint16 = 0xA11A             // Vocord Telecom
	WaveFormatVocordG726              uint16 = 0xA11B             // Vocord Telecom
	WaveFormatVocordG7221             uint16 = 0xA11C             // Vocord Telecom
	WaveFormatVocordG728              uint16 = 0xA11D             // Vocord Telecom
	WaveFormatVocordG729              uint16 = 0xA11E             // Vocord Telecom
	WaveFormatVocordG729A             uint16 = 0xA11F             // Vocord Telecom
	WaveFormatVocordG7231             uint16 = 0xA120             // Vocord Telecom
	WaveFormatVocordLbc               uint16 = 0xA121             // Vocord Telecom
	WaveFormatNiceG728                uint16 = 0xA122             // Nice Systems
	WaveFormatFraceTelecomG729        uint16 = 0xA123             // France Telecom
	WaveFormatCodian                  uint16 = 0xA124             // CODIAN
	WaveFormatFlac                    uint16 = 0xF1AC             // flac.sourceforge.net
	WaveFormatExtensible              uint16 = 0xFFFE             // Microsoft
	WaveFormatDevelopment             uint16 = 0xFFFF             // Development / Unregistered
)

func EncodeWave(samples pcm.Samples, waveFormat uint16) {
	context := samples.Context()

	fmtChunk := chunk{DwFourCC: chunkFmt}
	var fmtChunkData = bytes.Buffer{}

	switch waveFormat {
	case WaveFormatPcm, WaveFormatIeeeFloat, WaveFormatAlaw, WaveFormatMulaw:
		fmtData := PcmWaveFormat{
			Wf: WaveFormat{
				WFormatTag:      waveFormat,
				NChannels:       uint16(context.Channels),
				NSamplesPerSec:  uint32(context.SampleRate),
				NAvgBytesPerSec: uint32(context.SampleRate * context.Channels * samples.BitPerSample() / 8),
				NBlockAlign:     uint16(context.Channels * samples.BitPerSample() / 8),
			},
			WBitsPerSample: uint16(samples.BitPerSample()),
		}
		_ = binary.Write(&fmtChunkData, binary.LittleEndian, fmtData)
	}

	fmtChunk.DwSize = uint32(fmtChunkData.Len())

	dataChunk := chunk{DwFourCC: chunkData}
	var dataChunkData = bytes.Buffer{}

	switch v := samples.Pcm().(type) {
	case []uint8:
		for _, sample := range v {
			if err := binary.Write(&dataChunkData, binary.LittleEndian, sample); err != nil {
				fmt.Println(err)
			}
		}

	case []int16:
		for _, sample := range v {
			if err := binary.Write(&dataChunkData, binary.LittleEndian, sample); err != nil {
				fmt.Println(err)
			}
		}

	case []float32:
		for _, sample := range v {
			if err := binary.Write(&dataChunkData, binary.LittleEndian, sample); err != nil {
				fmt.Println(err)
			}
		}
	}

	dataChunk.DwSize = uint32(dataChunkData.Len())

	riffList := list{DwList: riffId, DwFourCC: formatWave}
	var riffData = bytes.Buffer{}

	_ = binary.Write(&riffData, binary.LittleEndian, fmtChunk)
	_, _ = riffData.ReadFrom(&fmtChunkData)

	_ = binary.Write(&riffData, binary.LittleEndian, dataChunk)
	_, _ = riffData.ReadFrom(&dataChunkData)

	riffList.DwSize = uint32(riffData.Len())

	outFile, _ := os.Create("out.wav")
	_ = binary.Write(outFile, binary.LittleEndian, riffList)
	_ = binary.Write(outFile, binary.LittleEndian, riffData.Bytes())
}
