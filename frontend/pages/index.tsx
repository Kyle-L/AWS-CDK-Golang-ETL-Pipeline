import { Button, Container, Fade, Flex, Heading, Select, Spinner, Table, TableContainer, Tag, Tbody, Td, Text, Th, Thead, Tr } from '@chakra-ui/react'
import { useEffect, useState } from 'react'
import FilterModal from '../components/FilterModal';
import TransactionRow from '../components/TransactionModal';
import { Transaction, TransactionQueryResponse } from '../types/types'
import { getTransactions } from '../utilities/transaction';

export default function Home() {
  const [filter, setFilter] = useState({ day: '', month: '1', year: '2016', isFraud: 'true', pageSize: '100' });
  const [data, setData] = useState<TransactionQueryResponse | null>(null)
  const [allRows, setAllRows] = useState<Transaction[]>([])

  const loadMore = () => {
      getTransactions({
        ...filter,
        paginationToken: data?.paginationToken
      }, (data) => {
        setData(data)
        setAllRows(allRows => [...allRows, ...data.items])
      });
  }

  useEffect(() => {
    setData(null)
    getTransactions(filter, (data) => {
      setData(data)
      setAllRows(data.items)
    });
  }, [filter])

  return (
    <Container paddingTop={10} maxW="container.xl">
      <Heading as="h1" size="xl" marginBottom={10}>Bank Transactions</Heading>
      <Flex alignItems="center" justify={"right"}>
        <FilterModal filter={filter} onFilterChange={setFilter} />
      </Flex>
      <Flex alignItems="center" justify={"space-between"} marginBottom={5}>
        <Text fontSize='xl' fontWeight='bold'>Showing {allRows.length} results for { new Date(`${filter.year}-${filter.month}-${filter.day}`).toLocaleDateString('en-US', { month: 'long', year: 'numeric', day: (filter.day ? '2-digit' : undefined) }) }</Text>
        </Flex>
      <TableContainer marginTop={5}>
        {!data && <Fade in={true}><Flex alignItems="center" justify="center" h="100%">
          <Spinner size={'xl'} speed="0.65s" color="blue.500" />
        </Flex></Fade>
        }
        {data && allRows.length > 0 && (
          <Fade in={true}>
            <Table variant='simple' size='sm'>
            <Thead>
              <Tr>
                <Th>Account #</Th>
                <Th>Customer Id</Th>
                <Th>Transaction Time</Th>
                <Th>Merchant</Th>
                <Th>Transaction Type</Th>
                <Th>Is Fraud?</Th>
                <Th>Transaction Amount</Th>
              </Tr>
            </Thead>
            <Tbody>
              {
                allRows.map((item) => (
                  <TransactionRow key={item.id} transaction={item} />
                ))
              }
            </Tbody>
          </Table>
          </Fade>
          )
        }
        {data && allRows.length === 0 && (
          <Fade in={true}><Text textAlign='center' fontSize='xl' fontWeight='bold' color='gray.500'>No data found for the given filter</Text></Fade>
        )}
      </TableContainer>
      {allRows.length > 0 && data?.paginationToken && (
        <Flex justifyContent='center' paddingTop={10} paddingBottom={10}>
          <Button colorScheme='blue' size='sm' onClick={loadMore}>Load More</Button>
        </Flex>
      )}
    </Container>
  )
}
