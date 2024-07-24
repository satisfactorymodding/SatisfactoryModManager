package websocket

import (
	"log/slog"
	"net/http"

	"github.com/spf13/viper"
	engineio_types "github.com/zishang520/engine.io/types"
	"github.com/zishang520/socket.io/socket"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/ficsitcli"
)

func ListenAndServeWebsocket() {
	httpMux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    "localhost:" + viper.GetString("websocket-port"),
		Handler: httpMux,
	}

	options := &socket.ServerOptions{}
	options.SetCors(&engineio_types.Cors{
		Origin: true, // Allow any origin
	})
	io := socket.NewServer(nil, options)
	httpMux.Handle("/socket.io/", io.ServeHandler(nil))

	_ = io.On("connection", func(data ...any) {
		client := data[0].(*socket.Socket)
		_ = client.On("installedMods", func(datas ...any) {
			lockfile, err := ficsitcli.FicsitCLI.GetSelectedInstallLockfile()
			if err != nil {
				slog.Error("failed to get lockfile", slog.Any("error", err))
				return
			}
			if lockfile == nil {
				slog.Error("no lockfile found for websocket call", slog.Any("error", err))
				return
			}
			installedMods := make(map[string]string)
			for modReference, info := range lockfile.Mods {
				installedMods[modReference] = info.Version
			}
			_ = client.Emit("installedMods", installedMods)
		})
	})

	err := httpServer.ListenAndServe()
	if err != nil {
		slog.Error("failed to start websocket server", slog.Any("err", err))
	}
}
