using System;
using System.IO;
using Microsoft.Win32;

namespace SMLLoader.Scripts
{
    public static class Files
    {
        public static void OpenFile(string filter, Action<string> action)
        {
            var openFileDialog = new OpenFileDialog
            {
                    Filter = filter
            };

            if (openFileDialog.ShowDialog() != null)
            {
                if (!string.IsNullOrEmpty(openFileDialog.FileName))
                {
                    action(openFileDialog.FileName);
                }
            }
        }

        public static void DirectoryCopy(string sourceDirName, string destDirName, bool copySubDirs)
        {
            // Get the subdirectories for the specified directory.
            var dir = new DirectoryInfo(sourceDirName);

            if (!dir.Exists)
            {
                throw new DirectoryNotFoundException(
                                                     "Source directory does not exist or could not be found: "
                                                     + sourceDirName);
            }

            var dirs = dir.GetDirectories();
            // If the destination directory doesn't exist, create it.
            if (!Directory.Exists(destDirName))
            {
                Directory.CreateDirectory(destDirName);
            }

            // Get the files in the directory and copy them to the new location.
            var files = dir.GetFiles();
            foreach (var file in files)
            {
                var temppath = Path.Combine(destDirName, file.Name);
                file.CopyTo(temppath, false);
            }

            // If copying subdirectories, copy them and their contents to new location.
            if (copySubDirs)
            {
                foreach (var subdir in dirs)
                {
                    var temppath = Path.Combine(destDirName, subdir.Name);
                    DirectoryCopy(subdir.FullName, temppath, copySubDirs);
                }
            }
        }
    }
}
