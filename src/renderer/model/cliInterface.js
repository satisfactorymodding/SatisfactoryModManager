import { exec } from 'child_process'

// TODO: SMLauncher CLI will be installed locally
const CLI_PATH = '..\\CLI\\SatisfactoryModLauncherCLI.exe'

const RunCommand = function (command, ...args) {
  return new Promise((resolve, reject) => {
    let process = [CLI_PATH, command]
    args.forEach(arg => process.push(arg))
    exec(process.join(' '), {}, (error, stdout, stderr) => {
      if (error) { return reject(error) }
      if (stderr) { return reject(stderr) }
      return resolve(stdout.trim())
    })
  })
}

export default {
  RunCommand
}
