using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.IO;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;
using Microsoft.Win32;
using Newtonsoft.Json;
using SMLLoader.Scripts;
using System.Runtime.InteropServices;
using System.Reflection;
using Newtonsoft.Json.Linq;
using Ionic.Zip;

namespace SMLLoader {
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window {
        private readonly string _configLocation = Environment.CurrentDirectory + "\\config.cfg";
        private readonly string _disabledModLoaderVersionsLocation = Environment.CurrentDirectory + "\\versions";
        private readonly string _disabledModsLocation = Environment.CurrentDirectory + "\\mods";
        private readonly string _modLoaderDll = "xinput1_3.dll";

        private Dictionary<CheckBox, TextBlock> _mods;
        private Config _config;
        private bool _handle = true;

        public MainWindow() {
            InitializeComponent();

            _config = new Config();
            _mods = new Dictionary<CheckBox, TextBlock>();
            LoadConfig();
        }

        /// <summary>
        /// Launch the game
        /// </summary>
        private void LaunchButton_Click(object sender, RoutedEventArgs e) {
            LaunchGame();
        }

        /// <summary>
        /// Open the mods folder
        /// </summary>
        private void OpenModsFolderButton_Click(object sender, RoutedEventArgs e) {
            OpenModsFolder();
        }

        /// <summary>
        /// Opens a dialog menu to add a new mod
        /// </summary>
        private void AddNewModButton_Click(object sender, RoutedEventArgs e) {
            AddNewMod();
        }

        /// <summary>
        /// Locate the game
        /// </summary>
        private void OpenGameLocationButton_Click(object sender, RoutedEventArgs e) {
            OpenGameLocation();
        }

        /// <summary>
        /// Relist all available mods
        /// </summary>
        private void ReloadModsButton_Click(object sender, RoutedEventArgs e) {
            ReloadMods();
        }

        private void AddNewModLoaderVersionButton_Click(object sender, RoutedEventArgs e) {
            AddNewModLoaderVersion();
        }

        private void GithubButton_Click(object sender, RoutedEventArgs e) {
            Process.Start(@"https://github.com/satisfactorymodding/SatisfactoryModLoader");
        }

        private void DiscordButton_Click(object sender, RoutedEventArgs e) {
            Process.Start(@"https://discord.gg/surVAY9");
        }

        private void ModLoaderVersionDropdown_SelectionChanged(object sender, SelectionChangedEventArgs e) {
            if (_handle) {
                HandleVersionDropdown();
            }
            _handle = true;
        }

        private void ModItemChecked(object sender, RoutedEventArgs e) {
            var inline = _mods[(CheckBox)sender].Inlines.ElementAt(0);
            string mod = new TextRange(inline.ContentStart, inline.ContentEnd).Text.Replace(" ", string.Empty);
            // move into mods directory
            string path = $"{_disabledModsLocation}\\{mod}";

            if (!Directory.Exists(path)) {
                return;
            }

            Directory.Move(path, $"{_config.ModsLocation}\\{mod}");

            ReloadMods();
        }

        private void ModItemUnchecked(object sender, RoutedEventArgs e) {
            var inline = _mods[(CheckBox)sender].Inlines.ElementAt(0);
            string mod = new TextRange(inline.ContentStart, inline.ContentEnd).Text.Replace(" ", string.Empty);
            // move into disabled mods directory
            string path = $"{_config.ModsLocation}\\{mod}";

            if (!Directory.Exists(path)) {
                return;
            }

            Directory.Move(path, $"{_disabledModsLocation}\\{mod}");

            ReloadMods();
        }

        private void LaunchGame() {
            if (!File.Exists($"{_config.BaseLocation}\\{_modLoaderDll}") || string.IsNullOrEmpty(_config.Version) || ModLoaderVersionDropdown.Items.Count == 0) {
                MessageBox.Show("A ModLoader version has not been assigned");
                return;
            }

            try {
                Process game = Process.Start(_config.ExeLocation);
                Close();
            } catch (Exception exception) {
                MessageBox.Show(exception.Message);
            }
        }

        private void OpenModsFolder() {
            if (!Directory.Exists(_config.ModsLocation)) {
                Directory.CreateDirectory(_config.ModsLocation);
            }
            Process.Start("explorer.exe", _config.ModsLocation);
        }

        private void AddNewMod() {
            OpenFileDialog openFileDialog = new OpenFileDialog() {
                Multiselect = true,
                Filter = "(Zip Files)|*.zip"
            };

            try {
                if (openFileDialog.ShowDialog() != null) {
                    if (openFileDialog.FileNames.Length > 0) {
                        foreach (string name in openFileDialog.FileNames) {
                           // move over
                           // File.Move(name, _config.ModsLocation + System.IO.Path.GetFileName(name));
                           // extract zip
                           using(ZipFile zip = ZipFile.Read(name)) {
                                // extract the folder
                                zip.ExtractAll(_config.ModsLocation + System.IO.Path.GetFileNameWithoutExtension(name));
                            }
                        }
                        ReloadMods();
                    }
                }
            } catch (Exception exception) {
                MessageBox.Show(exception.Message);
            }
        }

        private void AddNewModLoaderVersion() {
            OpenFileDialog openFileDialog = new OpenFileDialog() {
                Filter = "(Dll Files)|*.dll"
            };

            try {
                if (openFileDialog.ShowDialog() != null) {
                    // move over
                    File.Move(openFileDialog.FileName, $"{_disabledModLoaderVersionsLocation}\\{System.IO.Path.GetFileName(openFileDialog.FileName)}");
                }
            } catch (Exception exception) {
                MessageBox.Show(exception.Message);
            }

            ReloadVersions();
            ReloadMods();
        }

        private void OpenGameLocation() {
            OpenFileDialog openFileDialog = new OpenFileDialog();

            try {
                if (openFileDialog.ShowDialog() != null) {
                    string fileName = openFileDialog.FileName;
                    _config.BaseLocation = System.IO.Path.GetDirectoryName(fileName);
                    _config.ExeLocation = _config.BaseLocation + "\\FactoryGame-Win64-Shipping.exe";
                    _config.ModsLocation = _config.BaseLocation + "\\mods\\";
                    ExeLocationTextBox.Text = _config.ExeLocation;
                    SaveConfig();
                }
            } catch (Exception exception) {
                MessageBox.Show(exception.Message);
            }
        }

        private void ReloadMods() {
            _mods.Clear();
            ModListComboBox.Items.Clear();
            if(string.IsNullOrEmpty(_config.ModsLocation)) {
                return;
            }

            var directories = Directory.GetDirectories(_config.ModsLocation);

            foreach (string directory in directories) {
                LoadMods(directory, true);
            }

            directories = Directory.GetDirectories(_disabledModsLocation);
            foreach (string directory in directories) {
                LoadMods(directory, false);
            }

            ModListComboBox.Items.SortDescriptions.Add(
                new System.ComponentModel.SortDescription("Title",
                System.ComponentModel.ListSortDirection.Ascending));
        }

        private void ReloadVersions() {
            var files = Directory.GetFiles(_disabledModLoaderVersionsLocation).Where(f => System.IO.Path.GetExtension(f) == ".dll");
            foreach (string path in files) {
                ModLoaderVersionDropdown.Items.Add(System.IO.Path.GetFileNameWithoutExtension(path));
            }

            ModLoaderVersionDropdown.SelectedIndex = 0;

            for (int i = 0; i < ModLoaderVersionDropdown.Items.Count; i++) {
                if ((string)ModLoaderVersionDropdown.Items[i] == _config.Version) {
                    ModLoaderVersionDropdown.SelectedItem = ModLoaderVersionDropdown.Items[i];
                    break;
                }
            }

            HandleVersionDropdown();
        }

        private void SaveConfig() {
            using (StreamWriter writer = new StreamWriter(_configLocation)) {
                writer.Write(JsonConvert.SerializeObject(_config, Formatting.Indented));
            }
        }

        private void LoadConfig() {
            if (!File.Exists(_configLocation)) {
                SaveConfig();
            }

            using (StreamReader reader = new StreamReader(_configLocation)) {
                _config = JsonConvert.DeserializeObject<Config>(reader.ReadToEnd());
            }

            ExeLocationTextBox.Text = _config.ExeLocation;

            // load in all mod versions
            if (!Directory.Exists(_disabledModLoaderVersionsLocation)) {
                Directory.CreateDirectory(_disabledModLoaderVersionsLocation);
            }

            ReloadVersions();
            ReloadMods();
        }

        private void HandleVersionDropdown() {
            if (ModLoaderVersionDropdown.Items.Count > 0) {
                string name = ModLoaderVersionDropdown.SelectedItem.ToString() + ".dll";
                string newName = "xinput1_3.dll";
                // copy over
                if (File.Exists($"{_disabledModLoaderVersionsLocation}\\{name}")) {
                    File.Copy($"{_disabledModLoaderVersionsLocation}\\{name}", $"{_config.BaseLocation}\\{newName}", true);
                }
                _config.Version = ModLoaderVersionDropdown.SelectedItem.ToString();
            } else {
                _config.Version = string.Empty;
            }
            SaveConfig();
            ReloadMods();

            Title = $"Satisfactory Mod Loader - {_config.Version}";
        }

        private void LoadMods(string directory, bool enabled) {
            string[] files = Directory.GetFiles(directory);

            JObject json = null;

            foreach (string path in files) {
                if (path.EndsWith(".cfg")) {
                    using (StreamReader reader = new StreamReader(path)) {
                        json = JObject.Parse(reader.ReadToEnd());
                    }
                    continue;
                } 

                if(!path.EndsWith(".dll")) {
                    continue;
                }

                Grid item = CreateModItem(json, directory, enabled, json["LauncherVersion"].Value<string>() == (_config.Version.Length == 0 ? string.Empty : _config.Version.Substring(1)));
                ModListComboBox.Items.Add(item);
            }
        }

        private static ImageSource BitmapFromUri(Uri source) {
            var bitmap = new BitmapImage();
            bitmap.BeginInit();
            bitmap.UriSource = source;
            bitmap.CacheOption = BitmapCacheOption.OnLoad;
            bitmap.EndInit();
            return bitmap;
        }

        private Grid CreateModItem(JObject json, string path, bool enabled, bool validMod) {
            // <Grid Height="55" Width="240" Background="White">
            Grid grid = new Grid() {
                Background = Brushes.Transparent,
                HorizontalAlignment = HorizontalAlignment.Stretch,
                Height = 45
            };
            CheckBox checkBox = new CheckBox() {
                Width = 20,
                Height = 20,
                HorizontalAlignment = HorizontalAlignment.Left,
                VerticalAlignment = VerticalAlignment.Center,
                Margin = new Thickness(3, 4, 0, 0),
                IsChecked = enabled
            };
            checkBox.Checked += ModItemChecked;
            checkBox.Unchecked += ModItemUnchecked;
            // <Image HorizontalAlignment="Left" Height="55" VerticalAlignment="Top" Width="55"/>
            ImageSource img = null;
            try {
                img = BitmapFromUri(new Uri($"{path}\\{json["Icon"].Value<string>()}"));
            } catch (Exception exception) {
                img = BitmapFromUri(new Uri("pack://application:,,,/SMLLoader;component/Images/oee.png"));
            }

            Image image = new Image() {
                Source = img,
                HorizontalAlignment = HorizontalAlignment.Left,
                VerticalAlignment = VerticalAlignment.Stretch,
                Width = 55,
                Margin = new Thickness(20, 0, 0, 0)
            };

            RenderOptions.SetBitmapScalingMode(image, BitmapScalingMode.HighQuality);
            // <Label Content="Label" VerticalContentAlignment="Center" HorizontalAlignment="Left" Margin="60,8,0,0" VerticalAlignment="Top" Height="35" Width="170" FontSize="20"/>
            TextBlock label = new TextBlock() {
                Name = "Title",
                FontSize = 16,
                Foreground = Brushes.White,
                HorizontalAlignment = HorizontalAlignment.Left,
                VerticalAlignment = VerticalAlignment.Top,
                TextAlignment = TextAlignment.Left,
                TextWrapping = TextWrapping.Wrap,
                Margin = new Thickness(70, 0, 0, 0),
                Width = 170,
                Height = 50
            };
            label.Inlines.Add(new Bold(new Run(json["Name"].Value<string>())));
            TextBlock versionLabel = new TextBlock() {
                Text = json["Version"].Value<string>(),
                Name = "Version",
                FontSize = 12,
                Foreground = Brushes.White,
                HorizontalAlignment = HorizontalAlignment.Right,
                VerticalAlignment = VerticalAlignment.Top,
                TextAlignment = TextAlignment.Right,
                TextWrapping = TextWrapping.Wrap,
                Margin = new Thickness(0, 25, 0, 0),
                Width = 170,
                Height = 50
            };
            TextBlock loaderVersionLabel = new TextBlock() {
                Text = json["LauncherVersion"].Value<string>(),
                Name = "Version",
                FontSize = 12,
                Foreground = validMod ? Brushes.White : Brushes.Red,
                HorizontalAlignment = HorizontalAlignment.Right,
                VerticalAlignment = VerticalAlignment.Top,
                TextAlignment = TextAlignment.Right,
                TextWrapping = TextWrapping.Wrap,
                Margin = new Thickness(0, 0, 0, 0),
                Width = 170,
                Height = 50
            };
            TextBlock authorLabel = new TextBlock() {
                Text = json["Authors"].Value<string>(),
                Name = "Author",
                FontSize = 12,
                Foreground = Brushes.White,
                HorizontalAlignment = HorizontalAlignment.Left,
                VerticalAlignment = VerticalAlignment.Top,
                TextAlignment = TextAlignment.Left,
                TextWrapping = TextWrapping.Wrap,
                Margin = new Thickness(70, 25, 0, 0),
                Width = 300,
                Height = 50
            };

            grid.Children.Add(checkBox);
            grid.Children.Add(image);
            grid.Children.Add(label);
            grid.Children.Add(versionLabel);
            grid.Children.Add(loaderVersionLabel);
            grid.Children.Add(authorLabel);

            _mods.Add(checkBox, label);

            return grid;
        }
    }
}