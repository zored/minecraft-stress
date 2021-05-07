<?php

file_put_contents('requests', file_get_contents('php://input'), FILE_APPEND);
file_put_contents('requests', json_encode($_REQUEST), FILE_APPEND);
file_put_contents('requests', json_encode($_SERVER), FILE_APPEND);
file_put_contents('requests', json_encode($_POST), FILE_APPEND);
file_put_contents('requests', json_encode($_GET), FILE_APPEND);
