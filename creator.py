import pathlib as pl

paths = (
    pl.Path("data/1289/2466/2018/02/15/21/56"),
    pl.Path("data/1289/2466/2018/08/16/21/56"),
    pl.Path("data/1289/2466/2018/05/23/21/56"),
    pl.Path("data/1289/2466/2018/05/24/21/56"),
    pl.Path("data/1289/2466/2018/08/22/21/56"),
    pl.Path("data/1289/2466/2018/02/15/21/56"),
    pl.Path("data/3574/8644/2017/01/15/09/04"),
    pl.Path("data/3574/8644/2018/07/22/09/04"),
    pl.Path("data/3574/8644/2018/07/23/09/04"),
    pl.Path("data/2137/7123/2017/01/15/09/04"),
    pl.Path("data/2137/7123/2018/06/22/09/04"),
    pl.Path("data/2137/7123/2018/06/23/09/04"),
)

for path in paths:
    path.mkdir(parents=True, exist_ok=True)
