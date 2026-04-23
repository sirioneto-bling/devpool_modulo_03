<?php

namespace App\Client;

class NominatimClient extends BaseClient
{
    protected string $baseUrl = 'https://nominatim.openstreetmap.org';
    protected int $timeout = 15;
    protected array $defaultHeaders = [
        'User-Agent' => 'projeto-mvc-php/1.0',
        'Accept' => 'application/json',
    ];

    private array $defaultParams = [
        'format' => 'json',
        'addressdetails' => 1,
        'limit' => 10,
    ];

    /**
     * Busca endereços por texto livre.
     * Ex: "Rua Augusta, São Paulo"
     */
    public function search(string $query, int $limit = 10): array
    {
        $params = array_merge($this->defaultParams, [
            'q' => $query,
            'limit' => $limit,
        ]);

        return $this->get('/search', $params);
    }

    /**
     * Geocodificação reversa: coordenadas → endereço.
     */
    public function reverse(float $lat, float $lon, int $zoom = 18): array
    {
        $params = array_merge($this->defaultParams, [
            'lat' => $lat,
            'lon' => $lon,
            'zoom' => $zoom,
        ]);

        unset($params['limit']);

        return $this->get('/reverse', $params);
    }

    /**
     * Busca detalhes de um local pelo OSM ID.
     * Ex: osm_type = "R" (relation), osm_id = 123456
     */
    public function lookup(array $osmIds): array
    {
        $params = array_merge($this->defaultParams, [
            'osm_ids' => implode(',', $osmIds),
        ]);

        unset($params['limit']);

        return $this->get('/lookup', $params);
    }
}
