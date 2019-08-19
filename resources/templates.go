package resources

// func loadTemplates(resourcesPath string) (*template.Template, error) {
// 	tmplsPath := filepath.Join(resourcesPath, "templates")
// 	fis, err := ioutil.ReadDir(tmplsPath)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to read templates directory")
// 	}

// 	var paths []string
// 	for _, fi := range fis {
// 		if fi.IsDir() {
// 			continue
// 		}
// 		paths = append(paths, filepath.Join(tmplsPath, fi.Name()))
// 	}

// 	tmpls, err := template.New("").ParseFiles(paths...)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to parse templates")
// 	}
// 	return tmpls, nil
// }
