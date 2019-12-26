package gogurt

// meson is a python package.
// I suggest that we use pip install --prefix InstallDir(Meson{}) --cache-dir meson.
// Or maybe --target (nah, too weird)
// Crap... Can we still use url?

// Better:
// We create virtualenv environment in InstallDir
// We could then call install on untarred cache directory instead.
// Can then call 'meson' from it once path is set.

//  Could we call 'pip wheel --use-binary :all: <meson-unpack>' to build it? (minimise install)
// This would require having 'wheel', which is probably not a thing.
// BUT: we need to install virtualenv package anyway, so may as well :P

// Python 3.6 comes with virtualenv inbuilt.: python3 -m venv <dir>
// So, fix up Python, then everything falls into place.
