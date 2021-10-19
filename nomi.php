<?php
function sendError($code, $content){
    http_response_code($code);
    $error->error = true;
    $error->msg = $content;
    echo json_encode($error);
}

header("Content-Type: application/json; charset:utf-8");
//accept only get method
if ($_SERVER['REQUEST_METHOD'] === "GET"){
    //switch based on the type (name, names, exist)
    switch($_GET["type"]){
        case "names":
                //check if n exist and is not null 
                if (isset($_GET["n"])){
                    //che the bounds of n
                    if (intval($_GET["n"]) >= 2 && intval($_GET["n"])<= 100 ){
                        //read the file
                        $fileContent = file_get_contents("nomi.txt");
                        //remove all the \r\n in the file
                        $names = preg_split("/\r\n|\n|\r/", $fileContent);

                        $startWith = $_GET["start"];
                        $filteredNames = array();
                        if (isset($startWith)){
                            //iterate thru the names array and filter by substring (startWith)
                            $filteredNames =  array_filter($names, function ($name) use ($startWith) {
                                return strpos($name, $startWith) === 0;
                            });
                            //leave only the array's values (idk why but with some values array_filter will return a 
                            //object and not just an array), if you want to see what i mean remove the line below and pass
                            //as start parameter "abb"
                            $filteredNames = array_values($filteredNames);
                        }else{
                            //copy by reference names (i could copy it and dump names but whatever)
                            $filteredNames = &$names;
                        }

                        $randNames = array();
                        //iterate thru the num of n (url query param)
                        for ($i = 0; $i < $_GET["n"]; $i++){
                            //get a random number between 0 and length of name
                            //run the loop again if "check" is already in randNames
                            do{
                                $check= $filteredNames[rand()%count($filteredNames)];
                            }while(in_array($check, $randNames));
                            //append check to randNames
                            $randNames[]=$check; 
                        }

                        $resp-> error = false;
                        $resp-> content = $randNames;
                        echo json_encode($resp);
                    }else{
                        //range not satisfiable (416)
                        http_response_code(416);
                        $error->error=true;
                        $error->msg="n out of bounds";
                        $error->minval=2;
                        $error->maxval=100;
                        echo json_encode($error);
                    }
                }else{
                    //length required
                    http_response_code(411);
                    $error->error=true;
                    $error->msg="must assing n";
                    $error->minval=2;
                    $error->maxval=100;
                    echo json_encode($error);
                }
            break;
        case "name":
            $fileContent = file_get_contents("nomi.txt");
            $names = preg_split("/\r\n|\n|\r/", $fileContent);

            $startWith = $_GET["start"];
            $filteredNames = array();
            if (isset($startWith)){
                $filteredNames =  array_filter($names, function ($name) use ($startWith) {
                    return strpos($name, $startWith) === 0;
                });
                $filteredNames = array_values($filteredNames);
            }else{
                $filteredNames = &$names;
            }

            $resp-> error = false;
            $resp-> content = $filteredNames[rand()%count($filteredNames)];
            echo json_encode($resp);
            //documentation of this code is the same as names one so check that one out
            break;
        case "exist":
            $fileContent = file_get_contents("nomi.txt");
            $names = preg_split("/\r\n|\n|\r/", $fileContent);
            if (isset($_GET["toSearch"])){
                $resp->error = false;
                $resp->found = in_array($_GET["toSearch"], $names);
                echo json_encode($resp);
            }else{
                sendError(400,"'toSearch' must be defined");
            }
            break;
        default:
            sendError(400,"'type' must be names, name, exist");
    }
}else{
    sendError(405, "must use get method");
}
?>