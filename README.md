# listfilesinfold

A program (sorry about the name) to add padding on images (jpg-, png-) whos long side happen to exceed its short side by a certain factor.
For each such image, the program creates new padded image, with a suffix added to file-name.

2022-01-25: have tested it on Windows 10 only
2023-02-04: have added functionality to accept command line argument, "flag", -ratiolimit, of type float, that if ommited defaults to 2.0, also changed the background color of padd-area to white. Added some info printed when running program.