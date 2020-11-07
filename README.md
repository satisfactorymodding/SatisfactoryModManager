# Satisfactory Mod Manager (SMM)

Handles all the steps of installing mods, including download of SML and Bootstrapper. For questions, you can ask in the [Satisfactory Modding discord](https://discord.gg/TShj39G)

#### Installation

Download the latest release from https://github.com/satisfactorymodding/SatisfactoryModManager/releases

#### Usage

Just install it, choose an updated mod, and click install. Everything else is handled by SMM.

It has been reported by some users on Linux with Steam that SMM will download mods, but they will not show up in game.  A fix for this seems to be to update the launch options in steam

* Open Steam
* Library -> Satisfactory
* Right click Satisfactory -> Properties
* "SET LAUNCH OPTIONS..." button
* Enter in `WINEDLLOVERRIDES="msdia140.dll,xinput1_3.dll=n,b" %command%` -> OK
* Now Launch Satisfactory

#### Developmemnt
``` bash
# install dependencies
yarn install

# serve with hot reload at localhost:9080
yarn dev

# build electron application for production
yarn dist

# lint all JS/Vue component files in `src/`
yarn lint

```
