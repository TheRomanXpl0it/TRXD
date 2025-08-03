import React, { useEffect, useState } from 'react'
import { Check, ChevronsUpDown } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ScrollArea, ScrollBar } from '@/components/ui/scroll-area'
import {
  Command,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from '@/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { cn } from '@/lib/utils'
import { fetchCountries } from '@/lib/backend-interaction'

interface CountryProps {
  name: string
  iso3: string
  iso2: string
  emoji: string
}

interface LocationSelectorProps {
  disabled?: boolean
  onCountryChange?: (country: CountryProps | null) => void
}

const LocationSelector = ({ disabled, onCountryChange }: LocationSelectorProps) => {
  const [selectedCountry, setSelectedCountry] = useState<CountryProps | null>(null)
  const [openCountryDropdown, setOpenCountryDropdown] = useState(false)
  const [countriesData, setCountriesData] = useState<CountryProps[]>([])

  useEffect(() => {
    async function loadCountries() {
      const result = await fetchCountries() as unknown as CountryProps[]
      if (!Array.isArray(result)) {
        console.error("Unexpected response:", result)
        return
      }
      setCountriesData(result)
    }
    loadCountries()
  }, [])

  const handleCountrySelect = (country: CountryProps | null) => {
    setSelectedCountry(country)
    onCountryChange?.(country)
  }

  return (
    <div className="flex gap-4">
      <Popover open={openCountryDropdown} onOpenChange={setOpenCountryDropdown}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={openCountryDropdown}
            disabled={disabled}
            className="w-full justify-between"
          >
            {selectedCountry ? (
              <div className="flex items-center gap-2">
                <span>{selectedCountry.emoji}</span>
                <span>{selectedCountry.name}</span>
              </div>
            ) : (
              <span>Select Country...</span>
            )}
            <ChevronsUpDown className="h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="p-0">
          <Command>
            <CommandInput placeholder="Search country..." />
            <CommandList>
              <CommandEmpty>No country found.</CommandEmpty>
              <CommandGroup>
                <ScrollArea className="h-[300px]">
                  {countriesData.map((country) => (
                    <CommandItem
                      key={country.name}
                      value={country.name}
                      onSelect={() => {
                        handleCountrySelect(country)
                        setOpenCountryDropdown(false)
                      }}
                      className="flex cursor-pointer items-center justify-between text-sm"
                    >
                      <div className="flex items-center gap-2">
                        <span>{country.emoji}</span>
                        <span>{country.name}</span>
                      </div>
                      <Check
                        className={cn(
                          'h-4 w-4',
                          selectedCountry?.name === country.name
                            ? 'opacity-100'
                            : 'opacity-0'
                        )}
                      />
                    </CommandItem>
                  ))}
                  <ScrollBar orientation="vertical" />
                </ScrollArea>
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  )
}

export default LocationSelector
