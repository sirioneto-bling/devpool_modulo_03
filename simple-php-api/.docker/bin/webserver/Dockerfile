FROM php:8.2-apache-bookworm
ENV TZ=America/Sao_Paulo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ARG DEBIAN_FRONTEND=noninteractive

# Atualiza pacotes primeiro
RUN apt-get -y update && apt-get -y upgrade

# Instala ferramentas úteis
RUN apt-get -yq install apt-utils nano wget dialog

# Instala bibliotecas necessárias
RUN apt-get -y install build-essential git curl zip openssl libmemcached-dev zlib1g-dev

# Instala extensões PHP
RUN docker-php-ext-install pdo pdo_mysql bcmath opcache

# Composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

# Instala xdebug
RUN pecl install xdebug
RUN docker-php-ext-enable xdebug

# Timezone
RUN pecl install timezonedb
RUN docker-php-ext-enable timezonedb

# Locales
RUN apt-get install -y locales locales-all
RUN sed -i '/pt_BR.UTF-8/s/^# //g' /etc/locale.gen && locale-gen

# Limpa cache do apt
RUN apt-get clean && rm -rf /var/lib/apt/lists/*