package main

import (
	"encoding/json"
	"net/http"
	"time"
	"bytes"
	"fmt"
	"strings"
	"io/ioutil"
)


const template  = `
<file in='%d' %s interlace='' aspect_ratio='' scale_type='default' format_3d='default' right_eye_video='' external_audio='' external_audio_offset='' img_stub='' img.default_duration='3600.0' img.tone='' no_breaks='false' open_file.max_wait='0' open_url.max_wait='10000' force_mpeg_tc='false' audio_track='0' video_track='0' ts_program='-1' in_preroll='1.0' tc_preroll='1.0' network.buffer_min='1.0' network.buffer_max='10.0' network.buffer_wait_kf='false' network.max_video_gap='80' file.buffer_min='0.0' file.buffer_max='0.0' file.buffer_wait_kf='false' decode='auto' file.max_forward_rate='4.0' duration='%d' ref_id='%d' path='%s'>
        <info ts_programs='0' video_tracks='1' audio_tracks='1' streams='2' format='mov,mp4,m4a,3gp,3g2,mj2' format_name='QuickTime / MOV' start_time='0.04' duration='1746.861' size='561637955' bitrate='2572101' kbps_avg='319378.662109' kbps_avg_video='312159.667969' kbps_avg_audio='7218.994141' kbps_avg_data='0.0' metadata.major_brand='isom' metadata.minor_version='512' metadata.compatible_brands='isomiso2avc1mp41' metadata.encoder='Lavf54.29.104' network='false' buffer_video='-0.04' buffer_audio='0.25542'>
            <video.0 idx='0' codec='h264' codec_name='H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10' codec_tag='avc1' profile='High' width='1920' height='1080' has_b_frames='1' pixel_ar='1:1' display_ar='16:9' r_frame_rate='25/1' avg_frame_rate='25/1' time_base='1/25' start_pts='1' start_time='0.04' duration_ts='43670' duration='%d' bit_rate='2440263' nb_frames='43670' metadata.language='und' metadata.handler_name='VideoHandler'/>
            <audio.0 m_audio_track='0' idx='1' codec='aac' codec_name='AAC (Advanced Audio Coding)' codec_tag='mp4a' format='fltp' sample_rate='44100' channels='2' bits='32' time_base='1/44100' start_pts='1764' start_time='0.04' duration_ts='77036544' duration='%d' bit_rate='125591' nb_frames='75231' metadata.language='und' metadata.handler_name='SoundHandler'/>
        </info>
        <video width='1920' height='1080' rate='25.0' pixel_format='UYVY' aspect_ratio='16:9' interlace='progressive'/>
        <audio channels='2' rate='44100' bits='32' track_split='2'>
            <track source_index='0' mode='enabled' input_channels='0, 1' gain='0.0, 0.0' mute='0, 0' output_channels='0, 1' desc='Track 0'/>
        </audio>
        <object default_name='MFile' default_tracks='0' channels_per_track='0' internal.convert_frame='false' pause.fields='1' external_process='true' scaling_quality='8' crop='' mirror='' overlay_rms='false' overlay_rms.pos='0.1' overlay_rms.color='green' overlay_waveform='false' overlay_waveform.pos='-0.3' overlay_waveform.color='' mdelay.enabled='false' mdelay.live_preview='false'/>
        <mitem_props stop_in='' stop_out='' pause_in='' pause_out='' schedule_waitstart=''>
            <transition_in type='' time=''/>
            <transition_out type='' time=''/>
        </mitem_props>
    </file>`

type SmilPlaylist struct {
	Id               []byte `storm:"id"`
	Title            string
	Thumbnail        string
	Duration         string
	StartTimeSeconds int
	EndTimeSeconds   int
	VidType          videoType
	Scheduled        time.Time
	EndTime          time.Time
	Src              string
	StartSec         int
	Lenght           int
}

func GetPlaylist() ([]SmilPlaylist, error) {

	resp, err := http.Get("http://opinion.azorestv.com/api/fakelive/getplaylist")
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	var playlist []SmilPlaylist
	err = json.NewDecoder(resp.Body).Decode(&playlist)
	if err != nil {
		return nil, nil
	}

	return playlist, nil
}

func GenPlaylist(playlist []SmilPlaylist)error  {

	var ctr int

	buf := bytes.NewBuffer([]byte(""))

	buf.WriteString(`<m_config in_preroll='1.0' tc_preroll='1.0' path='C:\Users\Paulo Feliciano\Desktop\playlist\playlist.xml'>
    <video width='1920' height='1080' rate='25.0' pixel_format='UYVY' aspect_ratio='16:9' interlace='progressive'/>
    <audio channels='2' rate='44100' bits='32' track_split='2'>
        <track source_index='0' mode='enabled' input_channels='0, 1' gain='0.0, 0.0' mute='0, 0' output_channels='0, 1' desc='Track 0'/>
    </audio>`)

	for i:=range playlist{


		basedir := `V:\`

		file:=strings.Replace(playlist[i].Src,"/",`\`,-1)

		dur :=int(playlist[i].EndTime.Sub(playlist[i].Scheduled))/int(time.Second)

		var out=""
		if playlist[i].EndTimeSeconds!=0{
			out="out='"+fmt.Sprint(playlist[i].EndTimeSeconds) +"'"
		}



ctr++

		buf.WriteString(fmt.Sprintf(template,playlist[i].StartTimeSeconds,out,dur,ctr,basedir+file,dur,dur))

	}

	buf.WriteString(`</m_config>`)



	ioutil.WriteFile("playlist.xml",buf.Bytes(),0644)


	return nil
}
