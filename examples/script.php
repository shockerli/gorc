#!/usr/bin/env php
<?php

(new Script)->start();

class Script
{

    /**
     * @var bool ignore scalar type
     */
    public $ignoreType = false;

    /**
     * @var string log file path
     */
    public $logFile = __DIR__ . '/data.log';

    public function start()
    {
        $this->log('starting...');

        while (!feof(STDIN)) {
            // read input
            $input = fgets(STDIN);

            // decode data
            $data                              = json_decode($input, true);
            $data['request']['body']           = json_decode($data['request']['body'], true);
            $data['original_response']['body'] = $original = json_decode($data['original_response']['body'], true);
            $data['replayed_response']['body'] = $replayed = json_decode($data['replayed_response']['body'], true);

            // compare response body
            $this->log($this->sortArray($original) == $this->sortArray($replayed) ? 'Compare pass' : 'Compare fail');
            $this->log($data);
        }
    }

    private function log($data)
    {
        $data = is_scalar($data) ? $data : json_encode($data, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_PRETTY_PRINT);
        $data = date('Y-m-d H:i:s') . ' ' . $data;
        $data .= "\n";
        file_put_contents($this->logFile, $data, FILE_APPEND);
    }

    private function sortArray(array $data)
    {
        $data = array_map(function ($item) {
            switch (gettype($item)) {
                case 'array':
                    return $this->sortArray($item); // recursive array
                default:
                    return $this->ignoreType ? (string)$item : $item;
            }
        }, $data);

        // asc sort by key
        asort($data);

        return $data;
    }
}
