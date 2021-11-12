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
        $this->log('gorc script starting...');

        while (!feof(STDIN)) {
            // read input
            $input = fgets(STDIN);

            // decode data
            $data                              = json_decode($input, true);
            $data['request']['body']           = json_decode($data['request']['body'], true);
            $data['original_response']['body'] = $original = json_decode($data['original_response']['body'], true);
            $data['replayed_response']['body'] = $replayed = json_decode($data['replayed_response']['body'], true);

            // compare response body
            $compare = $this->sortArray($original) == $this->sortArray($replayed);
            $this->log(sprintf('[ReqId:%s] Compare %s', $data['req_id'], $compare ? 'pass' : 'fail'));
            $this->log($data);
        }
    }

    private function log($data)
    {
        $data = is_scalar($data) ? $data : json_encode($data, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);
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
