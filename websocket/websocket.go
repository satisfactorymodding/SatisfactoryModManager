//go:build !bindings

package websocket

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	engineio_types "github.com/zishang520/engine.io/types"
	"github.com/zishang520/socket.io/socket"

	"github.com/satisfactorymodding/SatisfactoryModManager/bindings"
)

func ListenAndServeWebsocket() {
	httpServer := engineio_types.CreateServer(nil)
	options := &socket.ServerOptions{}
	options.SetCors(&engineio_types.Cors{
		Origin: true, // Allow any origin
	})
	io := socket.NewServer(nil, options)
	httpServer.Handle("/socket.io", io.ServeHandler(nil))

	_ = io.On("connection", func(data ...any) {
		client := data[0].(*socket.Socket)
		_ = client.On("installedMods", func(datas ...any) {
			lockfile, err := bindings.BindingsInstance.FicsitCLI.GetLockFile(bindings.BindingsInstance.FicsitCLI.GetSelectedInstall())
			if err != nil {
				log.Error().Err(err).Msg("Failed to get lockfile")
				return
			}
			if lockfile == nil {
				log.Error().Err(err).Msg("no lockfile found for websocket call")
				return
			}
			installedMods := make(map[string]string)
			for modReference, info := range lockfile {
				installedMods[modReference] = info.Version
			}
			_ = client.Emit("installedMods", installedMods)
		})
	})

	httpServer.Listen("localhost:"+viper.GetString("websocket-port"), nil)
}
