import React, { useState, useMemo, useEffect, useCallback } from 'react';
import PropTypes from 'prop-types';
import { useTable, useGlobalFilter, useSortBy, useAsyncDebounce } from 'react-table';
import { useSelector, useDispatch } from 'react-redux';


// TODO: move this file closer to HostsDataTable
import hostInterface from 'interfaces/host';
import { humanHostMemory, humanHostUptime } from 'kolide/helpers';
import { setPagination } from 'redux/nodes/components/ManageHostsPage/actions';
import scrollToTop from 'utilities/scroll_to_top';
import InputField from 'components/forms/fields/InputField';
import HostPagination from 'components/hosts/HostPagination';

import TextCell from '../TextCell/TextCell';
import StatusCell from '../StatusCell/StatusCell';
import LinkCell from '../LinkCell/LinkCell';


// TODO: pull out to another file
// How we are handling lables and host counts on the client is strange. This function is required
// to try to hide some of that complexity, but ideally we'd come back and simplify how we are
// working with labels on the client.
const calculateTotalHostCount = (selectedFilter, labels, statusLabels) => {
  if (Object.keys(labels).length === 0) return 0;

  let hostCount = 0;
  switch (selectedFilter) {
    case 'all-hosts':
      hostCount = statusLabels.total_count;
      break;
    case 'new':
      hostCount = statusLabels.new_count;
      break;
    case 'online':
      hostCount = statusLabels.online_count;
      break;
    case 'offline':
      hostCount = statusLabels.offline_count;
      break;
    case 'mia':
      hostCount = statusLabels.mia_count;
      break;
    default: {
      const labelId = selectedFilter.split('/')[1];
      hostCount = labels[labelId].count;
      break;
    }
  }
  return hostCount;
};

// This data table uses react-table for implementation. The relevant documentation of the library
// can be found here https://react-table.tanstack.com/docs/api/useTable
const HostsDataTable = (props) => {
  console.log('render HostsDataTable');
  const { selectedFilter = '' } = props;
  const dispatch = useDispatch();
  const loadingHosts = useSelector(state => state.entities.hosts.loading);
  const hosts = useSelector(state => state.entities.hosts.data);
  const page = useSelector(state => state.components.ManageHostsPage.page);
  const perPage = useSelector(state => state.components.ManageHostsPage.perPage);
  const totalHostCount = useSelector((state) => {
    return calculateTotalHostCount(
      selectedFilter,
      state.entities.labels.data,
      state.components.ManageHostsPage.status_labels,
    );
  });

  useEffect(() => {
    console.log('fetching');
    dispatch(setPagination(page, perPage, selectedFilter));
  }, [dispatch, selectedFilter, page, perPage]);

  const columns = useMemo(() => {
    return [
      { Header: 'Hostname', accessor: 'hostname', Cell: cellProps => <LinkCell value={cellProps.cell.value} host={cellProps.row.original} /> },
      { Header: 'Status', accessor: 'status', Cell: cellProps => <StatusCell value={cellProps.cell.value} /> },
      { Header: 'OS', accessor: 'os_version', Cell: cellProps => <TextCell value={cellProps.cell.value} /> },
      { Header: 'Osquery', accessor: 'osquery_version', Cell: cellProps => <TextCell value={cellProps.cell.value} /> },
      { Header: 'IPv4', accessor: 'primary_ip', Cell: cellProps => <TextCell value={cellProps.cell.value} /> },
      { Header: 'Physical Address', accessor: 'primary_mac', Cell: cellProps => <TextCell value={cellProps.cell.value} /> },
      { Header: 'CPU', accessor: 'host_cpu', Cell: cellProps => <TextCell value={cellProps.cell.value} /> },
      { Header: 'Memory', accessor: 'memory', Cell: cellProps => <TextCell value={cellProps.cell.value} formatter={humanHostMemory} /> },
      { Header: 'Uptime', accessor: 'uptime', Cell: cellProps => <TextCell value={cellProps.cell.value} formatter={humanHostUptime} /> },
    ];
  }, []);

  const data = useMemo(() => {
    return Object.values(hosts);
  }, [hosts]);

  const {
    headerGroups,
    rows,
    prepareRow,
    setGlobalFilter,
    state,
  } = useTable({ columns, data }, useGlobalFilter, useSortBy);

  const [searchQuery, setSearchQuery] = useState('');

  const debouncedGlobalFilter = useAsyncDebounce((value) => {
    setGlobalFilter(value || undefined);
  }, 200);

  const onChange = useCallback((value) => {
    setSearchQuery(value);
    debouncedGlobalFilter(value);
  }, [setSearchQuery, debouncedGlobalFilter]);

  const onPaginationChange = useCallback((nextPage) => {
    dispatch(setPagination(nextPage, perPage, selectedFilter));
    scrollToTop();
  }, [dispatch, perPage, selectedFilter]);

  if (loadingHosts) return null;

  return (
    <React.Fragment>
      <InputField
        placeholder="Search hosts by hostname"
        name=""
        onChange={onChange}
        value={searchQuery}
        inputWrapperClass={'host-side-panel__filter-labels'}
      />

      {/* TODO: pull out into component */}
      <div className={'hosts-table hosts-table__wrapper'}>
        <table className={'hosts-table__table'}>
          <thead>
            {headerGroups.map(headerGroup => (
              <tr {...headerGroup.getHeaderGroupProps()}>
                {headerGroup.headers.map(column => (
                  <th {...column.getHeaderProps(column.getSortByToggleProps())}>
                    {column.render('Header')}
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {rows.map((row) => {
              prepareRow(row);
              return (
                <tr {...row.getRowProps()}>
                  {row.cells.map((cell) => {
                    return (
                      <td {...cell.getCellProps()}>
                        {cell.render('Cell')}
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      <HostPagination
        allHostCount={totalHostCount}
        currentPage={page}
        hostsPerPage={perPage}
        onPaginationChange={onPaginationChange}
      />
    </React.Fragment>
  );
};

HostsDataTable.propTypes = {
  hosts: PropTypes.arrayOf(hostInterface),
};

export default HostsDataTable;
