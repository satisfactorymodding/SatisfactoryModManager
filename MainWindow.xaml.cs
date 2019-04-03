using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Documents;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Threading;
using Ionic.Zip;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using SMLLoader.Scripts;

/*
 * #66252525
 */

namespace SMLLoader
{
    /// <summary>
    ///     Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        #region Constructors

        public MainWindow()
        {
            InitializeComponent();

            _config = new Config();
            _mods = new Dictionary<CheckBox, TextBlock>();

            ExeLocationTextBox.Text = string.Empty;

            DirectoryCheck();

            LoadConfig();
            ReloadVersions();

            _timer = new DispatcherTimer();
            _timer.Tick += Tick;
            _timer.Interval = new TimeSpan(0, 0, 1);
            _timer.Start();
        }

        #endregion Constructors

        #region Variables

        private static readonly string _appData = Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData) +
                                                  "\\Satisfactory Mod Launcher";

        private readonly string _configLocation = _appData + "\\config.cfg";
        private readonly string _disabledModLoaderVersionsLocation = _appData + "\\versions";
        private readonly string _disabledModsLocation = _appData + "\\mods";
        private readonly string _modLoaderDll = "xinput1_3.dll";

        private readonly Dictionary<CheckBox, TextBlock> _mods;
        private Config _config;
        private bool _handle = true;
        private readonly DispatcherTimer _timer;

        #endregion Variables

        #region EventHandlers

        private void Tick(object sender, EventArgs e)
        {
            ReloadMods();
        }

        // ui events

        private void LaunchButton_Click(object sender, RoutedEventArgs e)
        {
            LaunchGame();
        }

        private void OpenModsFolderButton_Click(object sender, RoutedEventArgs e)
        {
            OpenModsFolder();
        }

        private void AddNewModButton_Click(object sender, RoutedEventArgs e)
        {
            AddNewMod();
        }

        private void OpenGameLocationButton_Click(object sender, RoutedEventArgs e)
        {
            OpenGameLocation();
        }

        private void ReloadModsButton_Click(object sender, RoutedEventArgs e)
        {
            ReloadMods(true);
        }

        private void AddNewModLoaderVersionButton_Click(object sender, RoutedEventArgs e)
        {
            AddNewModLoaderVersion();
        }

        private void GithubButton_Click(object sender, RoutedEventArgs e)
        {
            Process.Start(@"https://github.com/satisfactorymodding/SatisfactoryModLoader");
        }

        private void DiscordButton_Click(object sender, RoutedEventArgs e)
        {
            Process.Start(@"https://discord.gg/surVAY9");
        }

        private void ModLoaderVersionDropdown_SelectionChanged(object sender, SelectionChangedEventArgs e)
        {
            if (_handle)
            {
                HandleVersionDropdown();
            }

            _handle = true;
        }

        private void ModItemChecked(object sender, RoutedEventArgs e)
        {
            DirectoryCheck();
            var inline = _mods[(CheckBox)sender].Inlines.ElementAt(0);
            var mod = new TextRange(inline.ContentStart, inline.ContentEnd).Text.Replace(" ", string.Empty);
            // move into mods directory
            var path = $"{_disabledModsLocation}\\{mod}";
            var modPath = $"{_config.ModsLocation}\\{mod}";

            if (!Directory.Exists(path))
            {
                return;
            }

            Files.DirectoryCopy(path, modPath, true);
            Directory.Delete(path, true);

            ReloadMods();
        }

        private void ModItemUnchecked(object sender, RoutedEventArgs e)
        {
            DirectoryCheck();
            var inline = _mods[(CheckBox)sender].Inlines.ElementAt(0);
            var mod = new TextRange(inline.ContentStart, inline.ContentEnd).Text.Replace(" ", string.Empty);
            // move into disabled mods directory
            var path = $"{_config.ModsLocation}\\{mod}";
            var modPath = $"{_disabledModsLocation}\\{mod}";

            if (!Directory.Exists(path))
            {
                return;
            }

            Files.DirectoryCopy(path, modPath, true);
            Directory.Delete(path, true);

            ReloadMods();
        }

        #endregion EventHandlers

        #region Methods

        private void LaunchGame()
        {
            if (string.IsNullOrEmpty(_config.BaseLocation) || string.IsNullOrEmpty(_config.ExeLocation) ||
                !File.Exists($"{_config.BaseLocation}\\{_modLoaderDll}") || string.IsNullOrEmpty(_config.Version) ||
                ModLoaderVersionDropdown.Items.Count == 0)
            {
                MessageBox.Show("A ModLoader version has not been assigned");
                return;
            }

            try
            {
                var game = Process.Start(_config.ExeLocation);
                Close();
            }
            catch (Exception exception)
            {
                MessageBox.Show(exception.Message);
            }
        }

        private void OpenModsFolder()
        {
            if (string.IsNullOrEmpty(_config.BaseLocation) || string.IsNullOrEmpty(_config.ModsLocation))
            {
                MessageBox.Show("Game directory not set.");
                return;
            }

            if (!Directory.Exists(_config.ModsLocation))
            {
                Directory.CreateDirectory(_config.ModsLocation);
            }

            Process.Start("explorer.exe", _config.ModsLocation);
        }

        private void AddNewMod()
        {
            if (string.IsNullOrEmpty(_config.BaseLocation) || string.IsNullOrEmpty(_config.ModsLocation))
            {
                MessageBox.Show("Game directory not set.");
                return;
            }

            DirectoryCheck();

            Files.OpenFile("(Zip Files)|*.zip", fileName =>
            {
                using (var zip = ZipFile.Read(fileName))
                {
                    // extract the folder
                    zip.ExtractAll(_config.ModsLocation + Path.GetFileNameWithoutExtension(fileName));
                }

                ReloadMods();
            });
        }

        // add new version from directory
        private void AddNewModLoaderVersion()
        {
            if (string.IsNullOrEmpty(_config.BaseLocation) || string.IsNullOrEmpty(_config.ModsLocation))
            {
                MessageBox.Show("Game directory not set.");
                return;
            }

            DirectoryCheck();

            Files.OpenFile("(Dll Files)|*.dll", fileName =>
            {
                var file = $"{_disabledModLoaderVersionsLocation}\\{Path.GetFileName(fileName)}";
                if (File.Exists(fileName))
                {
                    File.Move(fileName, file);
                    ReloadVersions();
                }
            });
        }

        // add game location from directory
        private void OpenGameLocation()
        {
            Files.OpenFile("(Exe files)|*.exe", fileName =>
            {
                try
                {
                    _config.BaseLocation = Path.GetDirectoryName(fileName);
                    _config.ExeLocation = _config.BaseLocation + "\\FactoryGame-Win64-Shipping.exe";
                    _config.ModsLocation = _config.BaseLocation + "\\mods\\";
                    ExeLocationTextBox.Text = _config.ExeLocation;
                    SaveConfig();
                }
                catch (Exception e)
                {
                    MessageBox.Show(e.Message);
                }
            });
        }

        private void ReloadMods(bool manual = false)
        {
            DirectoryCheck();

            // reset
            _mods.Clear();
            ModListComboBox.Items.Clear();

            // error catch
            if (string.IsNullOrEmpty(_config.ModsLocation))
            {
                if (manual)
                {
                    MessageBox.Show("Game directory not set.");
                }

                return;
            }

            // get mods
            var directories = Directory.GetDirectories(_config.ModsLocation);

            foreach (var directory in directories)
            {
                LoadMods(directory, true);
            }

            directories = Directory.GetDirectories(_disabledModsLocation);
            foreach (var directory in directories)
            {
                LoadMods(directory, false);
            }

            // sort by title
            ModListComboBox.Items.SortDescriptions.Add(
                                                       new SortDescription("Title",
                                                                           ListSortDirection.Ascending));
        }

        private void ReloadVersions()
        {
            DirectoryCheck();

            try
            {
                var files = Directory.GetFiles(_disabledModLoaderVersionsLocation)
                                     .Where(f => Path.GetExtension(f) == ".dll");
                foreach (var path in files)
                {
                    ModLoaderVersionDropdown.Items.Add(Path.GetFileNameWithoutExtension(path));
                }

                for (var i = 0; i < ModLoaderVersionDropdown.Items.Count; i++)
                {
                    if ((string)ModLoaderVersionDropdown.Items[i] == _config.Version)
                    {
                        ModLoaderVersionDropdown.SelectedIndex = i;
                        ModLoaderVersionDropdown.SelectedItem = ModLoaderVersionDropdown.Items[i];
                        break;
                    }
                }

                if (ModLoaderVersionDropdown.Items.Count > 0 && ModLoaderVersionDropdown.SelectedItem == null)
                {
                    ModLoaderVersionDropdown.SelectedIndex = 0;
                }

                HandleVersionDropdown();
            }
            catch (Exception e)
            {
                MessageBox.Show(e.Message);
            }
        }

        // utility

        private void SaveConfig()
        {
            using (var writer = new StreamWriter(_configLocation))
            {
                writer.Write(JsonConvert.SerializeObject(_config, Formatting.Indented));
            }
        }

        private void LoadConfig()
        {
            if (!File.Exists(_configLocation))
            {
                SaveConfig();
                return;
            }

            using (var reader = new StreamReader(_configLocation))
            {
                _config = JsonConvert.DeserializeObject<Config>(reader.ReadToEnd());
            }

            ExeLocationTextBox.Text = _config.ExeLocation;
        }

        private void HandleVersionDropdown()
        {
            if (ModLoaderVersionDropdown.Items.Count > 0)
            {
                var name = ModLoaderVersionDropdown.SelectedItem + ".dll";
                var newName = "xinput1_3.dll";
                var disabledPath = $"{_disabledModLoaderVersionsLocation}\\{name}";
                // copy over
                if (File.Exists(disabledPath))
                {
                    File.Copy(disabledPath, $"{_config.BaseLocation}\\{newName}", true);
                }

                _config.Version = ModLoaderVersionDropdown.SelectedItem.ToString();
            }
            else
            {
                _config.Version = string.Empty;
            }

            SaveConfig();
            ReloadMods();

            Title = $"Satisfactory Mod Loader - {_config.Version}";
        }

        private void LoadMods(string directory, bool enabled)
        {
            if (!Directory.Exists(directory))
            {
                MessageBox.Show("Invalid Directory");
                return;
            }

            var files = Directory.GetFiles(directory);

            JObject json = null;

            foreach (var path in files)
            {
                if (path.EndsWith(".cfg"))
                {
                    using (var reader = new StreamReader(path))
                    {
                        json = JObject.Parse(reader.ReadToEnd());
                    }

                    continue;
                }

                if (!path.EndsWith(".dll"))
                {
                    continue;
                }

                var item = CreateModItem(json, directory, enabled,
                                         json["LauncherVersion"].Value<string>() == (_config.Version.Length == 0
                                                 ? string.Empty
                                                 : _config.Version.Substring(1)));
                ModListComboBox.Items.Add(item);
            }
        }

        private static ImageSource BitmapFromUri(Uri source)
        {
            var bitmap = new BitmapImage();
            bitmap.BeginInit();
            bitmap.UriSource = source;
            bitmap.CacheOption = BitmapCacheOption.OnLoad;
            bitmap.EndInit();
            return bitmap;
        }

        private Grid CreateModItem(JObject json, string path, bool enabled, bool validMod)
        {
            // <Grid Height="55" Width="240" Background="White">
            var grid = new Grid
            {
                    Background = Brushes.Transparent,
                    HorizontalAlignment = HorizontalAlignment.Stretch,
                    Height = 45
            };
            var checkBox = new CheckBox
            {
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
            try
            {
                img = BitmapFromUri(new Uri($"{path}\\{json["Icon"].Value<string>()}"));
            }
            catch (Exception exception)
            {
                img = BitmapFromUri(new Uri("pack://application:,,,/SMLLoader;component/Images/oee.png"));
            }

            var image = new Image
            {
                    Source = img,
                    HorizontalAlignment = HorizontalAlignment.Left,
                    VerticalAlignment = VerticalAlignment.Stretch,
                    Width = 55,
                    Margin = new Thickness(20, 0, 0, 0)
            };

            RenderOptions.SetBitmapScalingMode(image, BitmapScalingMode.HighQuality);
            // <Label Content="Label" VerticalContentAlignment="Center" HorizontalAlignment="Left" Margin="60,8,0,0" VerticalAlignment="Top" Height="35" Width="170" FontSize="20"/>
            var label = new TextBlock
            {
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
            var versionLabel = new TextBlock
            {
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
            var loaderVersionLabel = new TextBlock
            {
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
            var authorLabel = new TextBlock
            {
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

        private void DirectoryCheck()
        {
            CreateDirectory(_appData);

            CreateDirectory(_disabledModLoaderVersionsLocation);
            CreateDirectory(_disabledModsLocation);

            if (!string.IsNullOrEmpty(_config.ModsLocation))
            {
                CreateDirectory(_config.ModsLocation);
            }
        }

        private void CreateDirectory(string path)
        {
            if (!Directory.Exists(path))
            {
                Directory.CreateDirectory(path);
            }
        }

        private void CreateFile(string path)
        {
            if (!File.Exists(path))
            {
                File.Create(path);
            }
        }

        #endregion Methods
    }
}
