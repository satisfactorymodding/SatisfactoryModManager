import { exec } from 'child_process'

const CLI_PATH = '.\\SatisfactoryModLauncherCLI.exe'

const RunCommand = function (command, ...args) {
  return new Promise((resolve, reject) => {
    const process = [CLI_PATH, command]
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
