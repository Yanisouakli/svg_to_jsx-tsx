package main


import (

  "fmt"
  "os"
  "path/filepath"
  "strings"
)

func extractSVGContent(svg string) string {
  start:= strings.Index(svg,"<svg")
  if start == -1{
      return svg
  }
  return svg[start:]
}


func generateDefaultJSX(fileName , svg string ) string{
    return fmt.Sprintf(`
import React from "react"
const %s = ({...props}) = (
%s
);
export default %s;
`,fileName, extractSVGContent(svg),fileName) 
}

func generateMUIJSX(fileName, svg string) string {
	return fmt.Sprintf(`import React from "react";
import { createSvgIcon } from "@mui/material";

const %s = createSvgIcon(
  %s,
  "%s"
);

export default %s;
`, fileName, extractSVGContent(svg), fileName, fileName)
}



func ProcessSvgFile(filePath,outputDir,outputType, fileType string) error { 
  content,err:= os.ReadFile(filePath)
  if err != nil {
    return fmt.Errorf("could not read file : %v", err)
  } 
svgContent:= string(content)
fileName := strings.TrimSuffix(filepath.Base(filePath),".svg")
output := ""

if outputType =="default"{
  output = generateDefaultJSX(fileName,svgContent)
} else if outputType =="mui" {
    output = generateMUIJSX(fileName,svgContent)
}

outputFilePath:= filepath.Join(outputDir,fmt.Sprintf("%s.%s",fileName,fileType)) 

err = os.WriteFile(outputFilePath, []byte(output),06444)
if err!= nil {
  return fmt.Errorf("could not WriteFile : %v",err)
}
  return nil
}


func main(){
  
	const (
		reset  = "\033[0m"    
		red    = "\033[31m"   
		green  = "\033[32m"   
    yellow = "\033[33m"   
		cyan   = "\033[36m"  
	)

  if len(os.Args) < 5 {
    fmt.Println(cyan + "Usage: " + reset + "go run main.go " + green + "<input-directory>" + reset + " " + green + "<output-directory>" + reset + " " + yellow + "<svg-type> " + yellow + "<file-type>" + reset )
  	fmt.Println(yellow + "<svg-type> " + reset + "can be " + green + "'default'" + reset + " or " + green + "'mui'" + reset)
    fmt.Println(yellow + "<file-type> " + reset + "can be " + green + "'jsx'" + reset + " or " + green + "'tsx'" + reset)
    os.Exit(1)
  }


  inputDir := os.Args[1]
  outputDir := os.Args[2] 
  outputType := os.Args[3]
  fileType := os.Args[4]

  if outputType != "default" && outputType != "mui"{
      fmt.Println(red + "Invalid output type. Use default or mui." + reset)
      os.Exit(1)
  }

  if fileType != "jsx" && fileType != "tsx" {
      fmt.Println(red + "Invalid file type. Use jsx or tsx. " + reset)
      os.Exit(1)
  }


  if _,err := os.Stat(outputDir); os.IsNotExist(err){
      err:= os.MkdirAll(outputDir,0755)
      if err!= nil {
        fmt.Printf("error creating output directory: %v\n",err)
        os.Exit(1)
      }

  }


  files,err:= os.ReadDir(inputDir)
  if err!= nil{
    fmt.Printf("oh oh error:{%v\n}",err)
    os.Exit(1)
  } 
   

  for _,file := range files {
     if strings.HasSuffix(file.Name(),".svg") {
       err:= ProcessSvgFile(filepath.Join(inputDir,file.Name()),outputDir,outputType,fileType)   
       if err != nil {
         fmt.Printf("Error processing the file %s: v%\n",file.Name(),err)

       } else {
         fmt.Printf("Processed file: %s\n",file.Name())
       }
     }
      
  }

}
