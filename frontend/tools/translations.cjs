const child_process = require('child_process');
const fs = require('fs');
const path = require('path');
require('dotenv').config();

const generatedDir = path.join(__dirname, '../src/lib/generated');
const i18nDir = path.join(generatedDir, 'i18n');
const i18nFile = path.join(i18nDir, 'index.ts');

if (fs.existsSync(i18nDir)) {
  console.log('Clearing i18n directory');
  fs.rmSync(i18nDir, { recursive: true });
  fs.mkdirSync(i18nDir);
}

console.log('Pulling translations');
child_process.execSync(`pnpm translations:pull -ak ${process.env.VITE_TOLGEE_API_KEY}`, { stdio: 'inherit' });

const langs = fs.readdirSync(i18nDir).map(file => file.replace('.json', ''));
console.log('Languages:', langs.join(', '));

const fileContent = '/* eslint-disable */\n' 
  + langs.map(lang => `import ${lang.replace('-', '_')} from './${lang}.json';`).join('\n') 
  + '\n\n' 
  + 'export const i18n = {\n' 
  + langs.map(lang => `  "${lang}": ${lang.replace('-', '_')},`).join('\n') 
  + '\n};\n';

fs.writeFileSync(i18nFile, fileContent);
console.log('Translations generated');
