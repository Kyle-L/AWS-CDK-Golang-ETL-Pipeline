import { EditIcon } from '@chakra-ui/icons';
import { Button, FormControl, FormHelperText, FormLabel, Heading, Modal, ModalBody, ModalCloseButton, ModalContent, ModalFooter, ModalHeader, ModalOverlay, Select, Stack, useDisclosure } from "@chakra-ui/react";
import { getDaysInMonth } from "../utilities/dateUtils";

export default function FilterModal({ filter, onFilterChange }: { filter: any, onFilterChange: any }) {
    const { isOpen, onOpen, onClose } = useDisclosure()
    const startYear = 2015;

    const handleSizeClick = () => {
        onOpen()
    }

    return (
        <>
            <Button size='sm' colorScheme={'gray'} leftIcon={<EditIcon />} onClick={handleSizeClick}>Apply Filters</Button>

            <Modal onClose={onClose} size={'xl'} isOpen={isOpen}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>
                        <Heading size="lg">Apply Filter</Heading>
                    </ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Stack spacing={3}>
                            <FormControl>
                                <FormLabel>Year</FormLabel>
                                <Select value={filter.year} onChange={(e) => onFilterChange({ ...filter, year: e.target.value })}>
                                    {Array.from(Array(new Date().getFullYear() + 1 - startYear).keys()).map((year) => (
                                        <option key={year} value={year + startYear}>{year + startYear}</option>
                                    ))}
                                </Select>
                            </FormControl>
                            <FormControl>
                                <FormLabel>Month</FormLabel>
                                <Select value={filter.month} onChange={(e) => onFilterChange({ ...filter, month: e.target.value })}>
                                    <option value='1'>January</option>
                                    <option value='2'>February</option>
                                    <option value='3'>March</option>
                                    <option value='4'>April</option>
                                    <option value='5'>May</option>
                                    <option value='6'>June</option>
                                    <option value='7'>July</option>
                                    <option value='8'>August</option>
                                    <option value='9'>September</option>
                                    <option value='10'>October</option>
                                    <option value='11'>November</option>
                                    <option value='12'>December</option>
                                </Select>
                            </FormControl>
                            <FormControl>
                                <FormLabel>Day of the Month</FormLabel>
                                <Select value={filter.day} onChange={(e) => onFilterChange({ ...filter, day: e.target.value })}>
                                    <option value=''>All Days</option>
                                    {getDaysInMonth(Number(filter.month), Number(filter.year)).map((day) => (
                                        <option key={day.getDate()} value={day.getDate()}>{day.getDate()}</option>
                                    ))}
                                </Select>
                            </FormControl>
                            <FormControl>
                                <FormLabel>Is Fraud</FormLabel>
                                <Select value={filter.isFraud} onChange={(e) => onFilterChange({ ...filter, isFraud: e.target.value })}>
                                    <option value='true'>Yes</option>
                                    <option value='false'>No</option>
                                </Select>
                            </FormControl>
                            <FormControl>
                                <FormLabel>Page Size</FormLabel>
                                <Select value={filter.pageSize} onChange={(e) => onFilterChange({ ...filter, pageSize: e.target.value })}>
                                    <option value='10'>10</option>
                                    <option value='20'>25</option>
                                    <option value='50'>50</option>
                                    <option value='100'>100</option>
                                    <option value='200'>250</option>
                                    <option value='200'>200</option>
                                    <option value='500'>250</option>
                                </Select>
                            </FormControl>
                        </Stack>
                    </ModalBody>
                    <ModalFooter>
                        <Button onClick={onClose}>Close</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </>
    )
}