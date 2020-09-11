const { spawn } = require('child_process');

const childProcess = spawn('yarn', [`dist:${process.platform.match(/[a-z]+/g)[0]}`], { shell: true, stdio: 'inherit' });
childProcess.on('exit', (code) => {
  process.exit(code);
});
