<?php

namespace App\Controllers\Api;

use App\Core\Controller;
use App\Client\NominatimClient;

class SearchAddressController extends Controller
{
    protected NominatimClient $nominatimClient;

    public function __construct()
    {
        $this->nominatimClient = new NominatimClient();
    }

    /**
     * GET /api/searchAddress/search?q=Rua+Augusta,+São+Paulo
     */
    public function search()
    {
        $this->validateRequestMethods(['GET']);

        $query = $_GET['q'] ?? null;
        $limit = (int) ($_GET['limit'] ?? 10);

        if (empty($query)) {
            return $this->jsonResponse(
                ['q' => 'O parâmetro "q" é obrigatório'],
                'error',
                400
            );
        }

        $result = $this->nominatimClient->search($query, $limit);

        if (!$result['success']) {
            return $this->jsonResponse(
                $result['error'] ?? 'Erro ao buscar endereço',
                'error',
                $result['statusCode'] ?: 502
            );
        }

        return $this->jsonResponse($result['data'], 'Busca realizada com sucesso');
    }

    /**
     * GET /api/searchAddress/reverse?lat=-23.56&lon=-46.65
     */
    public function reverse()
    {
        $this->validateRequestMethods(['GET']);

        $lat = $_GET['lat'] ?? null;
        $lon = $_GET['lon'] ?? null;

        if ($lat === null || $lon === null) {
            return $this->jsonResponse(
                ['lat' => 'obrigatório', 'lon' => 'obrigatório'],
                'error',
                400
            );
        }

        $result = $this->nominatimClient->reverse((float) $lat, (float) $lon);

        if (!$result['success']) {
            return $this->jsonResponse(
                $result['error'] ?? 'Erro na geocodificação reversa',
                'error',
                $result['statusCode'] ?: 502
            );
        }

        return $this->jsonResponse($result['data'], 'Geocodificação reversa realizada');
    }

    /**
     * GET /api/searchAddress/lookup?osm_ids=R146656,W104393803
     */
    public function lookup()
    {
        $this->validateRequestMethods(['GET']);

        $osmIds = $_GET['osm_ids'] ?? null;

        if (empty($osmIds)) {
            return $this->jsonResponse(
                ['osm_ids' => 'O parâmetro "osm_ids" é obrigatório (ex: R146656,W104393803)'],
                'error',
                400
            );
        }

        $ids = array_map('trim', explode(',', $osmIds));

        $result = $this->nominatimClient->lookup($ids);

        if (!$result['success']) {
            return $this->jsonResponse(
                $result['error'] ?? 'Erro ao buscar por OSM IDs',
                'error',
                $result['statusCode'] ?: 502
            );
        }

        return $this->jsonResponse($result['data'], 'Lookup realizado com sucesso');
    }
}
