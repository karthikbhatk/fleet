import endpoints from 'kolide/endpoints';

export default (client) => {
  return {
    destroy: (host) => {
      const { HOSTS } = endpoints;
      const endpoint = client._endpoint(`${HOSTS}/${host.id}`);

      return client.authenticatedDelete(endpoint);
    },
    load: (hostID) => {
      const { HOSTS } = endpoints;
      const endpoint = client._endpoint(`${HOSTS}/${hostID}`);

      return client.authenticatedGet(endpoint)
        .then(response => response.host);
    },

    // TODO: this needs to change to handle other filters
    loadAll: (page = 1, perPage = 100, selected = '', globalFilter = '', sortKey = '', isDesc) => {

      const { HOSTS, LABEL_HOSTS } = endpoints;

      const pagination = `page=${page - 1}&per_page=${perPage}`;

      let orderKey = '';
      if (sortKey !== '') {
        orderKey += `&order_key=${sortKey[0].id}`;
      }

      let orderDirection = '';
      if (isDesc !== undefined) {
        orderDirection = isDesc ? '&order_direction=desc' : '&order_direction=asc';
      }

      let searchQuery = '';
      if (globalFilter !== '') {
        searchQuery = `&query=${globalFilter}`;
      }

      let endpoint = '';
      const labelPrefix = 'labels/';
      if (selected.startsWith(labelPrefix)) {
        const lid = selected.substr(labelPrefix.length);
        endpoint = `${LABEL_HOSTS(lid)}?${pagination}`;
      } else {
        let selectedFilter = '';
        if (selected === 'new' || selected === 'online' || selected === 'offline' || selected === 'mia') {
          selectedFilter = `&status=${selected}`;
        }
        endpoint = `${HOSTS}?${pagination}${selectedFilter}${searchQuery}${orderKey}${orderDirection}`;
      }

      return client.authenticatedGet(client._endpoint(endpoint))
        .then((response) => { return response.hosts; });
    },
  };
};
