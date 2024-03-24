import Link from 'next/link'

export default function TestNav() {
    return (
        <nav>
            <ul>
                <li>
                    <Link href="/vehiclestatus">Home</Link>
                </li>
                <li>
                    <Link href="/vehiclestatus/123">123</Link>
                </li>
                <li>
                <Link href="/vehiclestatus/456">456</Link>
                </li>
            </ul>
        </nav>
    )
}