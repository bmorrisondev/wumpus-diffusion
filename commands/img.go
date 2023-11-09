package commands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

var ImgCommand CommandExt = CommandExt{
	Definition: discordgo.ApplicationCommand{
		Name:        "img",
		Description: "Generate an AI image",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "prompt",
				Description: "The Stable Diffusion prompt",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "height",
				Description: "Image height in pixels (default: 512)",
				Type:        discordgo.ApplicationCommandOptionInteger,
			},
			{
				Name:        "width",
				Description: "Image width in pixels (default: 512)",
				Type:        discordgo.ApplicationCommandOptionInteger,
			},
			{
				Name:        "negative",
				Description: "Negative prompt",
				Type:        discordgo.ApplicationCommandOptionString,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Got it! I'll drop a message back here when it's done.",
			},
		})

		b := Txt2ImgRequestBody{
			Prompt:     "",
			SendImages: true,
			SaveImages: true,
			Height:     512,
			Width:      512,
		}

		data := i.ApplicationCommandData()
		opts := data.Options

		for _, el := range opts {
			log.Printf("Option: %v | Value: %v", el.Name, el.Value)
			if el.Name == "prompt" {
				b.Prompt = el.StringValue()
			}
			if el.Name == "height" {
				b.Height = el.IntValue()
			}
			if el.Name == "width" {
				b.Width = el.IntValue()
			}
			if el.Name == "negative" {
				val := el.StringValue()
				b.NegativePrompt = &val
			}
		}
		jbytes, err := json.Marshal(b)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%v/sdapi/v1/txt2img", os.Getenv("SDW_BASE")), bytes.NewBuffer(jbytes))
		if err != nil {
			log.Fatal(err)
		}
		c := http.Client{}
		res, err := c.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()
		var rb Txt2ImgResponseBody
		err = json.NewDecoder(res.Body).Decode(&rb)
		if err != nil {
			log.Fatal(err)
		}

		decoded, err := base64.StdEncoding.DecodeString(rb.Images[0])
		if err != nil {
			log.Fatal(err)
		}

		ch, err := s.ThreadStart(i.ChannelID, fmt.Sprintf("Image request by @%v", i.Member.User.Username), discordgo.ChannelTypeGuildPublicThread, 1440)
		if err != nil {
			log.Fatal(err)
		}
		r := bytes.NewReader(decoded)
		s.ChannelMessageSend(ch.ID, fmt.Sprintf("<@%v>", i.Member.User.ID))
		s.ChannelFileSend(ch.ID, "generated.png", r)
		indented, err := json.MarshalIndent(b, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelMessageSend(ch.ID, fmt.Sprintf("```json\n%v\n```", string(indented)))
	},
}

type Txt2ImgRequestBody struct {
	Prompt         string  `json:"prompt"`
	SendImages     bool    `json:"send_images"`
	SaveImages     bool    `json:"save_images"`
	Height         int64   `json:"height"`
	Width          int64   `json:"width"`
	NegativePrompt *string `json:"negative"`
}

type Txt2ImgResponseBody struct {
	Images     []string `json:"images"`
	Parameters struct {
		Prompt                            string  `json:"prompt"`
		NegativePrompt                    string  `json:"negative_prompt"`
		Styles                            any     `json:"styles"`
		Seed                              int     `json:"seed"`
		Subseed                           int     `json:"subseed"`
		SubseedStrength                   int     `json:"subseed_strength"`
		SeedResizeFromH                   int     `json:"seed_resize_from_h"`
		SeedResizeFromW                   int     `json:"seed_resize_from_w"`
		SamplerName                       any     `json:"sampler_name"`
		BatchSize                         int     `json:"batch_size"`
		NIter                             int     `json:"n_iter"`
		Steps                             int     `json:"steps"`
		CfgScale                          float64 `json:"cfg_scale"`
		Width                             int     `json:"width"`
		Height                            int     `json:"height"`
		RestoreFaces                      any     `json:"restore_faces"`
		Tiling                            any     `json:"tiling"`
		DoNotSaveSamples                  bool    `json:"do_not_save_samples"`
		DoNotSaveGrid                     bool    `json:"do_not_save_grid"`
		Eta                               any     `json:"eta"`
		DenoisingStrength                 int     `json:"denoising_strength"`
		SMinUncond                        any     `json:"s_min_uncond"`
		SChurn                            any     `json:"s_churn"`
		STmax                             any     `json:"s_tmax"`
		STmin                             any     `json:"s_tmin"`
		SNoise                            any     `json:"s_noise"`
		OverrideSettings                  any     `json:"override_settings"`
		OverrideSettingsRestoreAfterwards bool    `json:"override_settings_restore_afterwards"`
		RefinerCheckpoint                 any     `json:"refiner_checkpoint"`
		RefinerSwitchAt                   any     `json:"refiner_switch_at"`
		DisableExtraNetworks              bool    `json:"disable_extra_networks"`
		Comments                          any     `json:"comments"`
		EnableHr                          bool    `json:"enable_hr"`
		FirstphaseWidth                   int     `json:"firstphase_width"`
		FirstphaseHeight                  int     `json:"firstphase_height"`
		HrScale                           float64 `json:"hr_scale"`
		HrUpscaler                        any     `json:"hr_upscaler"`
		HrSecondPassSteps                 int     `json:"hr_second_pass_steps"`
		HrResizeX                         int     `json:"hr_resize_x"`
		HrResizeY                         int     `json:"hr_resize_y"`
		HrCheckpointName                  any     `json:"hr_checkpoint_name"`
		HrSamplerName                     any     `json:"hr_sampler_name"`
		HrPrompt                          string  `json:"hr_prompt"`
		HrNegativePrompt                  string  `json:"hr_negative_prompt"`
		SamplerIndex                      string  `json:"sampler_index"`
		ScriptName                        any     `json:"script_name"`
		ScriptArgs                        []any   `json:"script_args"`
		SendImages                        bool    `json:"send_images"`
		SaveImages                        bool    `json:"save_images"`
		AlwaysonScripts                   struct {
		} `json:"alwayson_scripts"`
	} `json:"parameters"`
	Info string `json:"info"`
}
