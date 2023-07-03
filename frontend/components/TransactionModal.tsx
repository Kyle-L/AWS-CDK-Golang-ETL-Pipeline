import { Alert, AlertDescription, AlertIcon, AlertTitle, Button, Flex, Heading, Icon, Input, InputGroup, InputLeftAddon, InputRightAddon, Modal, ModalBody, ModalCloseButton, ModalContent, ModalFooter, ModalHeader, ModalOverlay, Select, Table, TableContainer, Tag, Tbody, Td, Text, Tr, useDisclosure } from "@chakra-ui/react"
import { useEffect, useRef, useState } from "react"
import { Transaction } from "../types/types"
import { debounce } from "lodash"
import { putTransaction } from "../utilities/transaction"
import countries from "countries-list";

export default function TransactionRow({ transaction }: { transaction: Transaction }) {
    const [transactionState, setTransactionState] = useState(transaction)
    const { isOpen, onOpen, onClose } = useDisclosure()

    const map: Record<string, any> = countries.countries;
    const mapKeys = Object.keys(map);

    const debouncedSave = useRef(debounce((transaction: Transaction) => {
        putTransaction(transaction, () => { });
    }, 500)).current;

    useEffect(() => {
        // Only update after changes have been made to the transaction.
        if (transactionState === transaction) {
            return;
        }

        // Applies a 500 ms debounce to the save so that we don't spam the API.
        debouncedSave(transactionState)
    }, [transaction, transactionState, debouncedSave])

    return (
        <>
            <Tr onClick={onOpen} sx={{ cursor: 'pointer', '&:hover': { backgroundColor: 'gray.100' } }}>
                <Td>{transactionState.accountNumber}</Td>
                <Td>{transactionState.customerId}</Td>
                <Td>{new Date(transactionState.transactionDateTime).toLocaleString('en-US', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })}</Td>
                <Td>{transactionState.merchantName}</Td>
                <Td>{transactionState.transactionType}</Td>
                <Td><Tag colorScheme={transactionState.isFraud === 'TRUE' ? 'red' : 'green'}>{transactionState.isFraud}</Tag></Td>
                <Td><Tag colorScheme='green'>{Number(transactionState.transactionAmount).toLocaleString('en-US', { style: 'currency', currency: 'USD' })}</Tag></Td>
                <Td><Tag colorScheme={transactionState.deleted ? 'red' : 'gray'}>{!transactionState.deleted ? 'Active' : 'Deleted'}</Tag></Td>
            </Tr>

            <Modal onClose={onClose} size={'full'} isOpen={isOpen}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>
                        <Heading size="lg">Transaction Details</Heading>
                        <Heading size={'sm'} color={'grey'}>{new Date(transactionState.transactionDateTime).toLocaleString('en-US', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })}</Heading>
                    </ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Alert
                            flexDirection='column'
                            alignItems='center'
                            justifyContent='center'
                            textAlign='center'
                            colorScheme={transactionState.isFraud === 'TRUE' ? 'red' : 'green'}
                            marginBottom={5}
                            borderRadius={10}
                        >
                            <AlertIcon />
                            <AlertTitle>{transactionState.isFraud === 'TRUE' ? 'Fraudulent Transaction' : 'Not Fraudulent'}</AlertTitle>
                            <AlertDescription maxWidth='sm'>
                                {transactionState.isFraud === 'TRUE' ? 'This transaction has been flagged as fraudulent.' : 'This transaction has not been flagged as fraudulent.'}
                                {' '} If that is not correct, please correct it below.
                            </AlertDescription>
                        </Alert>
                        <Table size={'sm'}>
                            <Tbody>
                                <Tr>
                                    <Td><Text fontWeight="bold">Transaction Date / Time</Text></Td>
                                    <Td><Input type="datetime-local" value={transactionState.transactionDateTime} onChange={(e) => setTransactionState({ ...transactionState, transactionDateTime: e.target.value })} /></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Account Number:</Text></Td>
                                    <Td><Input type="text" value={transactionState.accountNumber} onChange={(e) => setTransactionState({ ...transactionState, accountNumber: e.target.value })} /></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Customer ID:</Text></Td>
                                    <Td><Input type="text" value={transactionState.customerId} onChange={(e) => setTransactionState({ ...transactionState, customerId: e.target.value })} /></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Merchant Name:</Text></Td>
                                    <Td><Input type="text" value={transactionState.merchantName} onChange={(e) => setTransactionState({ ...transactionState, merchantName: e.target.value })} /></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Acquisition Country:</Text></Td>
                                    <Td>
                                        <Select value={transactionState.acqCountry} onChange={(e) => setTransactionState({ ...transactionState, acqCountry: e.target.value })}>
                                            {
                                                mapKeys.map((countryCode) => {
                                                    return <option key={countryCode} value={countryCode}>{map[countryCode].name}</option>
                                                })
                                            }
                                        </Select>
                                    </Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Merchant Country:</Text></Td>
                                    <Td>
                                        <Select value={transactionState.merchantCountryCode} onChange={(e) => setTransactionState({ ...transactionState, merchantCountryCode: e.target.value })}>
                                            {
                                                mapKeys.map((countryCode) => {
                                                    return <option key={countryCode} value={countryCode}>{map[countryCode].name}</option>
                                                })
                                            }
                                        </Select>
                                    </Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Merchant Category Code:</Text></Td>
                                    <Td><Input type="text" value={transactionState.merchantCategoryCode} onChange={(e) => setTransactionState({ ...transactionState, merchantCategoryCode: e.target.value.toLowerCase() })} /></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Account Open Date:</Text></Td>
                                    <Td><Input type="datetime-local" value={transactionState.accountOpenDate} onChange={(e) => setTransactionState({ ...transactionState, accountOpenDate: e.target.value })} /></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Date of Last Address Change:</Text></Td>
                                    <Td><Input type="datetime-local" value={transactionState.dateOfLastAddressChange} onChange={(e) => setTransactionState({ ...transactionState, dateOfLastAddressChange: e.target.value })} /></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Card Number:</Text></Td>
                                    <Td>
                                        <InputGroup>
                                            <InputLeftAddon>**** ****</InputLeftAddon>
                                            <Input type="text" value={transactionState.cardLast4Digits} onChange={(e) => setTransactionState({ ...transactionState, cardLast4Digits: +e.target.value })} />
                                            <InputLeftAddon>CVV</InputLeftAddon>
                                            <Input type="text" value={transactionState.cardCVV} onChange={(e) => setTransactionState({ ...transactionState, cardCVV: +e.target.value })} />
                                        </InputGroup>
                                    </Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Transaction Type:</Text></Td>
                                    <Td><Input type="text" value={transactionState.transactionType} onChange={(e) => setTransactionState({ ...transactionState, transactionType: e.target.value.toUpperCase() })} /></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Card Present:</Text></Td>
                                    <Td><Select value={transactionState.cardPresent} onChange={(e) => setTransactionState({ ...transactionState, cardPresent: e.target.value })}>
                                        <option value="TRUE">TRUE</option>
                                        <option value="FALSE">FALSE</option>
                                    </Select></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Is Fraudalent:</Text></Td>
                                    <Td><Select value={transactionState.isFraud} onChange={(e) => setTransactionState({ ...transactionState, isFraud: e.target.value })}>
                                        <option value="TRUE">TRUE</option>
                                        <option value="FALSE">FALSE</option>
                                    </Select></Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Credit Limit:</Text></Td>
                                    <Td>
                                        <InputGroup>
                                            <InputLeftAddon>$</InputLeftAddon>
                                            <Input type="number" value={transactionState.creditLimit} onChange={(e) => setTransactionState({ ...transactionState, creditLimit: +e.target.value })} />
                                        </InputGroup>
                                    </Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Available Money:</Text></Td>
                                    <Td><InputGroup>
                                        <InputLeftAddon>$</InputLeftAddon>
                                        <Input type="number" value={transactionState.availableMoney} onChange={(e) => setTransactionState({ ...transactionState, availableMoney: +e.target.value })} />
                                    </InputGroup>
                                    </Td>
                                </Tr>
                                <Tr>
                                    <Td><Text fontWeight="bold">Transaction Amount:</Text></Td>
                                    <Td><InputGroup>
                                        <InputLeftAddon>$</InputLeftAddon>
                                        <Input type="number" value={transactionState.transactionAmount} onChange={(e) => setTransactionState({ ...transactionState, transactionAmount: +e.target.value })} />
                                    </InputGroup>
                                    </Td>
                                </Tr>
                            </Tbody>
                        </Table>
                    </ModalBody>
                    <ModalFooter>
                        <Button onClick={onClose}>Close</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </>
    )
}